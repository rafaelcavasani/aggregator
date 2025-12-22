package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Document representa um documento genérico no Elasticsearch
type Document struct {
	ID        string                 `json:"id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ElasticsearchClient encapsula operações do Elasticsearch
type ElasticsearchClient struct {
	client *elasticsearch.Client
}

// NewElasticsearchClient cria uma nova instância do cliente
func NewElasticsearchClient(addresses []string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente elasticsearch: %w", err)
	}

	// Verificar conexão
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("erro na resposta do elasticsearch: %s", res.String())
	}

	fmt.Println("✅ Conectado ao Elasticsearch com sucesso!")

	return &ElasticsearchClient{client: client}, nil
}

// CreateIndex cria um novo índice no Elasticsearch
func (ec *ElasticsearchClient) CreateIndex(ctx context.Context, indexName string, mapping map[string]interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("erro ao codificar mapping: %w", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  &buf,
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return fmt.Errorf("erro ao criar índice: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// Índice já existe é aceitável
		if res.StatusCode == 400 {
			log.Printf("⚠️  Índice '%s' já existe\n", indexName)
			return nil
		}
		return fmt.Errorf("erro ao criar índice: %s", res.String())
	}

	log.Printf("✅ Índice '%s' criado com sucesso!\n", indexName)
	return nil
}

// IndexDocument insere um documento no índice
func (ec *ElasticsearchClient) IndexDocument(ctx context.Context, indexName string, docID string, document interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(document); err != nil {
		return fmt.Errorf("erro ao codificar documento: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       &buf,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return fmt.Errorf("erro ao indexar documento: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao indexar documento: %s", res.String())
	}

	log.Printf("✅ Documento '%s' inserido no índice '%s'\n", docID, indexName)
	return nil
}

// UpdateDocument atualiza um documento existente
func (ec *ElasticsearchClient) UpdateDocument(ctx context.Context, indexName string, docID string, updates map[string]interface{}) error {
	updateDoc := map[string]interface{}{
		"doc": updates,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updateDoc); err != nil {
		return fmt.Errorf("erro ao codificar update: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       &buf,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return fmt.Errorf("erro ao atualizar documento: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao atualizar documento: %s", res.String())
	}

	log.Printf("✅ Documento '%s' atualizado no índice '%s'\n", docID, indexName)
	return nil
}

// GetDocument busca um documento por ID
func (ec *ElasticsearchClient) GetDocument(ctx context.Context, indexName string, docID string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documento: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("erro ao buscar documento: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return result, nil
}

// SearchDocuments busca documentos usando query
func (ec *ElasticsearchClient) SearchDocuments(ctx context.Context, indexName string, query map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("erro ao codificar query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  &buf,
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documentos: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("erro ao buscar documentos: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Extrair hits
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		documents[i] = hit.(map[string]interface{})
	}

	return documents, nil
}

// DeleteDocument remove um documento
func (ec *ElasticsearchClient) DeleteDocument(ctx context.Context, indexName string, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, ec.client)
	if err != nil {
		return fmt.Errorf("erro ao deletar documento: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("erro ao deletar documento: %s", res.String())
	}

	log.Printf("✅ Documento '%s' deletado do índice '%s'\n", docID, indexName)
	return nil
}

// QueryRequest representa uma requisição genérica para o Elasticsearch
type QueryRequest struct {
	Operation  string                 `json:"operation"` // create_index, index, update, get, search, delete
	Index      string                 `json:"index"`
	DocumentID string                 `json:"document_id,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
}

// QueryResponse representa a resposta da API
type QueryResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

var esClient *ElasticsearchClient

// handleQuery processa requisições genéricas para o Elasticsearch
func handleQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(QueryResponse{
			Success: false,
			Error:   "Método não permitido. Use POST",
		})
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao decodificar JSON: %v", err),
		})
		return
	}

	ctx := context.Background()
	var response QueryResponse

	switch req.Operation {
	case "create_index":
		err := esClient.CreateIndex(ctx, req.Index, req.Body)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{Success: true, Message: fmt.Sprintf("Índice '%s' criado com sucesso", req.Index)}
		}

	case "index":
		err := esClient.IndexDocument(ctx, req.Index, req.DocumentID, req.Body)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{Success: true, Message: fmt.Sprintf("Documento '%s' inserido com sucesso", req.DocumentID)}
		}

	case "update":
		err := esClient.UpdateDocument(ctx, req.Index, req.DocumentID, req.Body)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{Success: true, Message: fmt.Sprintf("Documento '%s' atualizado com sucesso", req.DocumentID)}
		}

	case "get":
		doc, err := esClient.GetDocument(ctx, req.Index, req.DocumentID)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{Success: true, Message: "Documento encontrado", Data: doc}
		}

	case "search":
		results, err := esClient.SearchDocuments(ctx, req.Index, req.Body)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{
				Success: true,
				Message: fmt.Sprintf("Encontrados %d documentos", len(results)),
				Data:    map[string]interface{}{"hits": results, "total": len(results)},
			}
		}

	case "delete":
		err := esClient.DeleteDocument(ctx, req.Index, req.DocumentID)
		if err != nil {
			response = QueryResponse{Success: false, Error: err.Error()}
		} else {
			response = QueryResponse{Success: true, Message: fmt.Sprintf("Documento '%s' deletado com sucesso", req.DocumentID)}
		}

	default:
		response = QueryResponse{
			Success: false,
			Error:   fmt.Sprintf("Operação '%s' não suportada. Use: create_index, index, update, get, search, delete", req.Operation),
		}
	}

	json.NewEncoder(w).Encode(response)
}

// healthHandler verifica se o serviço está rodando
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "data-aggregator",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// formatarMoeda formata um valor float64 como moeda brasileira
func formatarMoeda(valor float64) string {
	// Separar parte inteira e decimal
	parteInteira := int64(valor)
	parteDecimal := int((valor - float64(parteInteira)) * 100)

	// Formatar parte inteira com separador de milhar
	strInteira := fmt.Sprintf("%d", parteInteira)
	var resultado strings.Builder

	for i, digit := range strInteira {
		if i > 0 && (len(strInteira)-i)%3 == 0 {
			resultado.WriteString(".")
		}
		resultado.WriteRune(digit)
	}

	return fmt.Sprintf("R$ %s,%02d", resultado.String(), parteDecimal)
}

// saldoClienteHandler retorna o saldo formatado de um cliente
func saldoClienteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Método não permitido. Use POST",
		})
		return
	}

	var req struct {
		CodigoCliente string `json:"codigo_cliente"`
		DataInicio    string `json:"data_inicio"`
		DataFim       string `json:"data_fim"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Erro ao decodificar JSON: %v", err),
		})
		return
	}

	// Construir query
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"codigo_cliente.keyword": req.CodigoCliente,
						},
					},
					{
						"range": map[string]interface{}{
							"data_vencimento": map[string]interface{}{
								"gte": req.DataInicio,
								"lte": req.DataFim,
							},
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"saldo_total": map[string]interface{}{
				"scripted_metric": map[string]interface{}{
					"init_script":    "state.saldo_total = 0.0",
					"map_script":     "double saldo = doc['valor_original'].value; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } state.saldo_total += saldo;",
					"combine_script": "return state.saldo_total",
					"reduce_script":  "double total = 0; for (s in states) { total += s; } return Math.round(total * 100.0) / 100.0",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Erro ao codificar query: %v", err),
		})
		return
	}

	// Executar busca
	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Erro ao executar busca: %v", err),
		})
		return
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Erro ao decodificar resposta: %v", err),
		})
		return
	}

	// Extrair saldo e formatar
	aggs := result["aggregations"].(map[string]interface{})
	saldoTotal := aggs["saldo_total"].(map[string]interface{})["value"].(float64)
	hits := result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)

	response := map[string]interface{}{
		"codigo_cliente": req.CodigoCliente,
		"periodo": map[string]string{
			"inicio": req.DataInicio,
			"fim":    req.DataFim,
		},
		"total_recebiveis": int(hits),
		"saldo_total":      saldoTotal,
		"saldo_formatado":  formatarMoeda(saldoTotal),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	// Conectar ao Elasticsearch
	var err error
	esClient, err = NewElasticsearchClient([]string{"http://localhost:9200"})
	if err != nil {
		log.Fatal(err)
	}

	// Configurar rotas HTTP
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/saldo-cliente", saldoClienteHandler)

	// Iniciar servidor HTTP
	port := ":8080"
	fmt.Printf("🚀 Servidor HTTP iniciado em http://localhost%s\n", port)
	fmt.Printf("📝 Endpoint de query: POST http://localhost%s/query\n", port)
	fmt.Printf("💚 Health check: GET http://localhost%s/health\n", port)
	fmt.Println("\n✅ Servidor pronto para receber requisições!")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
