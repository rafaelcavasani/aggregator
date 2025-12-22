# Teste de Performance - Elasticsearch

## Objetivo
Avaliar a performance de queries Elasticsearch no índice `ciclo_vida_recebivel` utilizando o cliente com maior volume de dados.

## Ambiente
- **Elasticsearch**: 8.11.0
- **Índice**: ciclo_vida_recebivel
- **Total de documentos no índice**: ~4.2 milhões
- **Cliente testado**: CLI-10008 (813.358 recebíveis)
- **Período de análise**: 2025-01-01 a 2026-12-31
- **Data do teste**: 21/12/2025

## Metodologia
Foram executadas 3 queries diferentes no cliente com maior volume de registros (CLI-10008):

1. **Contagem simples** (_count endpoint)
2. **Agregação com cálculo Painless** (saldo total)
3. **Busca com ordenação** (100 documentos)

---

## Resultados

### 1. Contagem de Recebíveis
**Query**: `get_receivables_count_by_customer.http`

```http
GET http://localhost:9200/ciclo_vida_recebivel/_count
{
  "query": {
    "term": {
      "codigo_cliente.keyword": "CLI-10008"
    }
  }
}
```

**Resultado**:
- ✅ Total de registros: **813.358**
- ⏱️ Tempo de resposta: **52.96ms**
- 📊 Performance: **Excelente** - Query extremamente eficiente para contagem

---

### 2. Cálculo de Saldo Total (Painless Script)
**Query**: `get_customer_balance_by_date.http`

```http
GET http://localhost:9200/ciclo_vida_recebivel/_search
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {"term": {"codigo_cliente.keyword": "CLI-10008"}},
        {"range": {"data_vencimento": {"gte": "2025-01-01", "lte": "2026-12-31"}}}
      ]
    }
  },
  "aggs": {
    "saldo_total": {
      "scripted_metric": {
        "init_script": "state.saldo_total = 0.0",
        "map_script": "double saldo = doc['valor_original'].value; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } state.saldo_total += saldo;",
        "combine_script": "return state.saldo_total",
        "reduce_script": "double total = 0; for (s in states) { total += s; } return Math.round(total * 100.0) / 100.0"
      }
    }
  }
}
```

**Resultado**:
- 💰 Saldo calculado: **R$ 337.781.136,35**
- ⏱️ Tempo de resposta: **10.086,64ms** (~10 segundos)
- 📊 Performance: **Aceitável** - Script Painless processou todos os documentos
- ⚠️ **Observação**: Query mais pesada devido ao cálculo em cada documento

---

### 3. Busca com Ordenação
**Query**: `get_receivables_by_customer_and_due_date.http`

```http
GET http://localhost:9200/ciclo_vida_recebivel/_search
{
  "query": {
    "bool": {
      "must": [
        {"term": {"codigo_cliente.keyword": "CLI-10008"}},
        {"range": {"data_vencimento": {"gte": "2025-01-01", "lte": "2026-12-31"}}}
      ]
    }
  },
  "sort": [{"data_vencimento": {"order": "asc"}}],
  "size": 100
}
```

**Resultado**:
- 📄 Total de hits: **10.000** (limite padrão do ES)
- 📄 Documentos retornados: **100**
- ⏱️ Tempo de resposta: **81.39ms**
- 📊 Performance: **Excelente** - Query rápida mesmo com ordenação

---

## Análise Comparativa

| Query | Operação | Tempo (ms) | Performance |
|-------|----------|------------|-------------|
| **Count** | Contagem simples | 52.96 | ⭐⭐⭐⭐⭐ Excelente |
| **Aggregation** | Cálculo Painless (813K docs) | 10,086.64 | ⭐⭐⭐ Aceitável |
| **Search** | Busca + Sort (100 docs) | 81.39 | ⭐⭐⭐⭐⭐ Excelente |

---

## Conclusões

### 🎯 Pontos Positivos
1. **Contagem**: Extremamente rápida (~50ms) mesmo com 813K documentos
2. **Busca paginada**: Muito eficiente com ordenação (81ms para 100 docs)
3. **Escalabilidade**: Sistema mantém boa performance com milhões de documentos

### ⚠️ Pontos de Atenção
1. **Scripts Painless**: Agregações com scripts complexos são custosas
   - 10 segundos para processar 813K documentos
   - Acessa `params._source` para iterar nested arrays
   - Recomendação: Usar para relatórios batch, não queries em tempo real

### 💡 Recomendações

#### Para Queries em Tempo Real (< 1s)
- ✅ Use `_count` para contagens
- ✅ Use `_search` com `size` limitado e ordenação simples
- ❌ Evite `scripted_metric` com volume alto

#### Para Relatórios/Analytics (> 1s aceitável)
- ✅ Use `scripted_metric` para cálculos complexos
- ✅ Considere cache ou materialização de resultados
- ✅ Execute em background ou com feedback de progresso

#### Otimizações Sugeridas
1. **Desnormalizar dados**: Armazenar `saldo_disponivel` calculado no índice
2. **Index refresh interval**: Ajustar para cargas de escrita
3. **Shard allocation**: Considerar re-indexing com mais shards se volume crescer
4. **Cache warming**: Pre-executar queries frequentes após startup

---

## Próximos Passos
- [ ] Testar performance com cliente mediano (~200K docs)
- [ ] Testar performance com cliente pequeno (~10K docs)
- [ ] Avaliar impacto de desnormalização do saldo
- [ ] Benchmark com diferentes configurações de shards
- [ ] Monitorar uso de memória heap durante agregações
