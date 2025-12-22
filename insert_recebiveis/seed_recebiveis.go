package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
)

// Recebivel representa um recebível
type Recebivel struct {
	IDRecebivel           string         `json:"id_recebivel"`
	CodigoCliente         string         `json:"codigo_cliente"`
	CodigoProduto         int            `json:"codigo_produto"`
	CodigoProdutoParceiro int            `json:"codigo_produto_parceiro"`
	Modalidade            int            `json:"modalidade"`
	ValorOriginal         float64        `json:"valor_original"`
	DataVencimento        string         `json:"data_vencimento"`
	Cancelamentos         []Cancelamento `json:"cancelamentos,omitempty"`
	Negociacoes           []Negociacao   `json:"negociacoes,omitempty"`
}

// Cancelamento representa um cancelamento
type Cancelamento struct {
	IDCancelamento   string  `json:"id_cancelamento"`
	DataCancelamento string  `json:"data_cancelamento"`
	ValorCancelado   float64 `json:"valor_cancelado"`
	Motivo           string  `json:"motivo"`
}

// Negociacao representa uma negociação
type Negociacao struct {
	IDNegociacao   string  `json:"id_negociacao"`
	DataNegociacao string  `json:"data_negociacao"`
	ValorNegociado float64 `json:"valor_negociado"`
}

func main() {
	// Inicializar seed aleatório
	rand.Seed(time.Now().UnixNano())

	// Configurar cliente Elasticsearch
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar cliente Elasticsearch: %s", err))
	}

	// Verificar conexão
	res, err := es.Info()
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar ao Elasticsearch: %s", err))
	}
	defer res.Body.Close()

	fmt.Println("✅ Conectado ao Elasticsearch")

	ctx := context.Background()
	totalRecebiveis := 10000000
	indexName := "ciclo_vida_recebivel"

	// Listas de clientes fictícios
	clientes := []string{
		"CLI-10001", "CLI-10002", "CLI-10003", "CLI-10004", "CLI-10005",
		"CLI-10006", "CLI-10007", "CLI-10008", "CLI-10009", "CLI-10010",
		"CLI-10011", "CLI-10012", "CLI-10013", "CLI-10014", "CLI-10015",
		"CLI-10016", "CLI-10017", "CLI-10018", "CLI-10019", "CLI-10020",
	}

	// Definir percentuais
	// 20% negociação parcial, 5% negociação total
	// 15% cancelamento parcial, 5% cancelamento total
	// Total: 45% com alguma operação, 55% sem operações

	fmt.Println("🚀 Iniciando inserção de recebíveis com goroutines e bulk indexer...")

	startTime := time.Now()

	// Configurar BulkIndexer para inserções em lote com backpressure
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        es,
		NumWorkers:    4,                // Reduzido para evitar TooManyRequests
		FlushBytes:    2e+6,             // Flush a cada 2MB (menor = mais frequente)
		FlushInterval: 10 * time.Second, // Flush mais frequente
		OnError: func(ctx context.Context, err error) {
			// Implementar retry com backoff exponencial
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "Too Many Requests") {
				fmt.Printf("⚠️  Rate limit atingido, aguardando 2s...\n")
				time.Sleep(2 * time.Second)
			}
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar BulkIndexer: %s", err))
	}

	// Contador atômico para sucesso e erro
	var (
		countSuccessful uint64
		countFailed     uint64
	)

	// WaitGroup para aguardar todas as goroutines
	var wg sync.WaitGroup

	// Canal para limitar número de goroutines simultâneas (backpressure)
	semaphore := make(chan struct{}, 50) // Máximo 50 goroutines simultâneas (reduzido para evitar 429)

	// Gerar e inserir recebíveis em paralelo
	for i := 0; i < totalRecebiveis; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Adquirir permissão

		go func(index int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Liberar permissão

			// Gerar recebível
			recebivel := gerarRecebivelConcorrente(clientes, index, totalRecebiveis)

			// Serializar para JSON
			body, err := json.Marshal(recebivel)
			if err != nil {
				atomic.AddUint64(&countFailed, 1)
				fmt.Printf("❌ Erro ao serializar recebível %d: %s\n", index+1, err)
				return
			}

			// Adicionar ao BulkIndexer com retry logic
			maxRetries := 3
			for attempt := 0; attempt < maxRetries; attempt++ {
				err = bi.Add(
					ctx,
					esutil.BulkIndexerItem{
						Action:     "index",
						DocumentID: recebivel.IDRecebivel,
						Body:       bytes.NewReader(body),
						OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
							atomic.AddUint64(&countSuccessful, 1)
							current := atomic.LoadUint64(&countSuccessful)
							if current%10000 == 0 {
								fmt.Printf("✅ Inseridos %d/%d recebíveis (%.2f%%)\n", current, totalRecebiveis, float64(current)/float64(totalRecebiveis)*100)
							}
						},
						OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
							atomic.AddUint64(&countFailed, 1)
							if err != nil {
								fmt.Printf("❌ Erro ao inserir recebível: %s\n", err)
							} else {
								if res.Error.Type == "es_rejected_execution_exception" {
									fmt.Printf("⚠️  ES Rejected (429) - recebível será reprocessado\n")
								} else {
									fmt.Printf("❌ Erro ao inserir recebível: %s: %s\n", res.Error.Type, res.Error.Reason)
								}
							}
						},
					},
				)
				if err == nil {
					break // Sucesso, sair do loop de retry
				}
				// Backoff exponencial: 100ms, 200ms, 400ms
				if attempt < maxRetries-1 {
					waitTime := time.Duration(100*(1<<uint(attempt))) * time.Millisecond
					fmt.Printf("⚠️  Retry %d/%d após %v devido a: %s\n", attempt+1, maxRetries, waitTime, err)
					time.Sleep(waitTime)
				}
			}
			if err != nil {
				atomic.AddUint64(&countFailed, 1)
				fmt.Printf("❌ Falha após %d tentativas: %s\n", maxRetries, err)
			}
		}(i)
	}

	// Aguardar todas as goroutines terminarem
	wg.Wait()

	// Fechar o BulkIndexer e processar itens restantes
	if err := bi.Close(ctx); err != nil {
		panic(fmt.Sprintf("Erro ao fechar BulkIndexer: %s", err))
	}

	// Estatísticas do BulkIndexer
	biStats := bi.Stats()

	duration := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🎉 Finalizado em %s!\n", duration)
	fmt.Printf("📊 Estatísticas:\n")
	fmt.Printf("   - Total de recebíveis: %d\n", totalRecebiveis)
	fmt.Printf("   - Sucesso: %d\n", countSuccessful)
	fmt.Printf("   - Falhas: %d\n", countFailed)
	fmt.Printf("   - BulkIndexer - Flushed: %d\n", biStats.NumFlushed)
	fmt.Printf("   - BulkIndexer - Indexed: %d\n", biStats.NumIndexed)
	fmt.Printf("   - BulkIndexer - Failed: %d\n", biStats.NumFailed)
	fmt.Printf("   - Taxa: %.0f recebíveis/segundo\n", float64(totalRecebiveis)/duration.Seconds())
	fmt.Println(strings.Repeat("=", 60))

	// Forçar refresh final
	es.Indices.Refresh(es.Indices.Refresh.WithIndex(indexName))
}

// gerarRecebivelConcorrente gera um recebível com dados aleatórios (thread-safe)
func gerarRecebivelConcorrente(clientes []string, index int, total int) Recebivel {
	// Criar gerador de números aleatórios específico para esta goroutine
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(index)))

	// Gerar UUID único
	id := uuid.New().String()

	// Selecionar cliente aleatório
	cliente := clientes[rng.Intn(len(clientes))]

	// Valor original entre 100 e 1000
	valorOriginal := float64(rng.Intn(901) + 100) // 100 a 1000

	// Gerar data de vencimento aleatória entre 2025-01-01 e 2026-12-31
	dataVencimento := gerarDataAleatoriaConcorrente(rng)

	recebivel := Recebivel{
		IDRecebivel:           id,
		CodigoCliente:         cliente,
		CodigoProduto:         rng.Intn(500) + 100,
		CodigoProdutoParceiro: rng.Intn(100) + 1,
		Modalidade:            rng.Intn(5) + 1,
		ValorOriginal:         valorOriginal,
		DataVencimento:        dataVencimento,
	}

	// Calcular percentual baseado no índice para distribuição uniforme
	percentual := float64(index) / float64(total) * 100

	// 5% cancelamento total (0-5%)
	if percentual < 5 {
		recebivel.Cancelamentos = gerarCancelamentoTotalConcorrente(valorOriginal, rng)
	} else if percentual < 20 { // 15% cancelamento parcial (5-20%)
		recebivel.Cancelamentos = gerarCancelamentoParcialConcorrente(valorOriginal, rng)
	} else if percentual < 25 { // 5% negociação total (20-25%)
		recebivel.Negociacoes = gerarNegociacaoTotalConcorrente(valorOriginal, rng)
	} else if percentual < 45 { // 20% negociação parcial (25-45%)
		recebivel.Negociacoes = gerarNegociacaoParcialConcorrente(valorOriginal, rng)
	}
	// 55% sem operações (45-100%)

	return recebivel
}

// gerarDataAleatoriaConcorrente gera uma data entre 2025-01-01 e 2026-12-31 (thread-safe)
func gerarDataAleatoriaConcorrente(rng *rand.Rand) string {
	inicio := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	delta := fim.Unix() - inicio.Unix()
	sec := rng.Int63n(delta)
	data := inicio.Add(time.Duration(sec) * time.Second)

	return data.Format("2006-01-02")
}

// gerarCancelamentoTotalConcorrente gera cancelamentos que somam 100% do valor original (thread-safe)
func gerarCancelamentoTotalConcorrente(valorOriginal float64, rng *rand.Rand) []Cancelamento {
	numCancelamentos := rng.Intn(3) + 1 // 1 a 3 cancelamentos
	cancelamentos := make([]Cancelamento, numCancelamentos)

	valorRestante := valorOriginal
	for i := 0; i < numCancelamentos; i++ {
		var valorCancelado float64
		if i == numCancelamentos-1 {
			// Último cancelamento pega o valor restante
			valorCancelado = valorRestante
		} else {
			// Cancelamentos intermediários pegam de 20% a 60% do restante
			percentual := rng.Float64()*0.4 + 0.2 // 0.2 a 0.6
			valorCancelado = valorRestante * percentual
			valorRestante -= valorCancelado
		}

		cancelamentos[i] = Cancelamento{
			IDCancelamento:   uuid.New().String(),
			DataCancelamento: gerarDataAleatoriaConcorrente(rng),
			ValorCancelado:   arredondar(valorCancelado),
			Motivo:           getMotivoAleatorioConcorrente(rng),
		}
	}

	return cancelamentos
}

// gerarCancelamentoParcialConcorrente gera cancelamentos que somam entre 10% e 70% do valor original (thread-safe)
func gerarCancelamentoParcialConcorrente(valorOriginal float64, rng *rand.Rand) []Cancelamento {
	numCancelamentos := rng.Intn(3) + 1 // 1 a 3 cancelamentos
	cancelamentos := make([]Cancelamento, numCancelamentos)

	// Total a cancelar: entre 10% e 70% do valor original
	percentualTotal := rng.Float64()*0.6 + 0.1 // 0.1 a 0.7
	totalACancelar := valorOriginal * percentualTotal

	valorRestante := totalACancelar
	for i := 0; i < numCancelamentos; i++ {
		var valorCancelado float64
		if i == numCancelamentos-1 {
			valorCancelado = valorRestante
		} else {
			percentual := rng.Float64()*0.4 + 0.2 // 0.2 a 0.6
			valorCancelado = valorRestante * percentual
			valorRestante -= valorCancelado
		}

		cancelamentos[i] = Cancelamento{
			IDCancelamento:   uuid.New().String(),
			DataCancelamento: gerarDataAleatoriaConcorrente(rng),
			ValorCancelado:   arredondar(valorCancelado),
			Motivo:           getMotivoAleatorioConcorrente(rng),
		}
	}

	return cancelamentos
}

// gerarNegociacaoTotalConcorrente gera negociações que somam 100% do valor original (thread-safe)
func gerarNegociacaoTotalConcorrente(valorOriginal float64, rng *rand.Rand) []Negociacao {
	numNegociacoes := rng.Intn(3) + 1 // 1 a 3 negociações
	negociacoes := make([]Negociacao, numNegociacoes)

	valorRestante := valorOriginal
	for i := 0; i < numNegociacoes; i++ {
		var valorNegociado float64
		if i == numNegociacoes-1 {
			valorNegociado = valorRestante
		} else {
			percentual := rng.Float64()*0.4 + 0.2 // 0.2 a 0.6
			valorNegociado = valorRestante * percentual
			valorRestante -= valorNegociado
		}

		negociacoes[i] = Negociacao{
			IDNegociacao:   uuid.New().String(),
			DataNegociacao: gerarDataAleatoriaConcorrente(rng),
			ValorNegociado: arredondar(valorNegociado),
		}
	}

	return negociacoes
}

// gerarNegociacaoParcialConcorrente gera negociações que somam entre 10% e 70% do valor original (thread-safe)
func gerarNegociacaoParcialConcorrente(valorOriginal float64, rng *rand.Rand) []Negociacao {
	numNegociacoes := rng.Intn(3) + 1 // 1 a 3 negociações
	negociacoes := make([]Negociacao, numNegociacoes)

	// Total a negociar: entre 10% e 70% do valor original
	percentualTotal := rng.Float64()*0.6 + 0.1 // 0.1 a 0.7
	totalANegociar := valorOriginal * percentualTotal

	valorRestante := totalANegociar
	for i := 0; i < numNegociacoes; i++ {
		var valorNegociado float64
		if i == numNegociacoes-1 {
			valorNegociado = valorRestante
		} else {
			percentual := rng.Float64()*0.4 + 0.2 // 0.2 a 0.6
			valorNegociado = valorRestante * percentual
			valorRestante -= valorNegociado
		}

		negociacoes[i] = Negociacao{
			IDNegociacao:   uuid.New().String(),
			DataNegociacao: gerarDataAleatoriaConcorrente(rng),
			ValorNegociado: arredondar(valorNegociado),
		}
	}

	return negociacoes
}

// arredondar arredonda para 2 casas decimais
func arredondar(valor float64) float64 {
	return float64(int(valor*100)) / 100
}

// getMotivoAleatorioConcorrente retorna um motivo aleatório para cancelamento (thread-safe)
func getMotivoAleatorioConcorrente(rng *rand.Rand) string {
	motivos := []string{
		"Cliente solicitou cancelamento parcial.",
		"Ajuste de valor por erro operacional.",
		"Negociação comercial com o cliente.",
		"Desconto promocional aplicado.",
		"Cancelamento por inadimplência.",
		"Renegociação de dívida.",
		"Ajuste contratual.",
	}
	return motivos[rng.Intn(len(motivos))]
}
