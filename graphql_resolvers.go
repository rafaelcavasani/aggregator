package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
)

// Resolver para buscar todos os recebíveis com limite
func getAllReceivablesResolver(params graphql.ResolveParams) (interface{}, error) {
	size, _ := params.Args["size"].(int)
	if size == 0 {
		size = 10
	}

	query := map[string]interface{}{
		"size": size,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	receivables := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		source["id"] = hitMap["_id"]
		receivables[i] = source
	}

	return map[string]interface{}{
		"total":       total,
		"receivables": receivables,
	}, nil
}

// Resolver para buscar recebível por ID
func getReceivableByIdResolver(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id é obrigatório")
	}

	doc, err := esClient.GetDocument(context.Background(), "ciclo_vida_recebivel", id)
	if err != nil {
		return nil, err
	}

	source := doc["_source"].(map[string]interface{})
	source["id"] = doc["_id"]
	return source, nil
}

// Resolver para buscar saldo do cliente por período
func getCustomerBalanceResolver(params graphql.ResolveParams) (interface{}, error) {
	codigoCliente, _ := params.Args["codigo_cliente"].(string)
	dataInicio, _ := params.Args["data_inicio"].(string)
	dataFim, _ := params.Args["data_fim"].(string)

	if codigoCliente == "" || dataInicio == "" || dataFim == "" {
		return nil, fmt.Errorf("codigo_cliente, data_inicio e data_fim são obrigatórios")
	}

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"codigo_cliente": codigoCliente,
						},
					},
					{
						"range": map[string]interface{}{
							"data_vencimento": map[string]interface{}{
								"gte": dataInicio,
								"lte": dataFim,
							},
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"resultado": map[string]interface{}{
				"filters": map[string]interface{}{
					"filters": map[string]interface{}{
						"all": map[string]interface{}{
							"match_all": map[string]interface{}{},
						},
					},
				},
				"aggs": map[string]interface{}{
					"soma_valores_originais": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "valor_original",
						},
					},
					"soma_cancelamentos": map[string]interface{}{
						"nested": map[string]interface{}{
							"path": "cancelamentos",
						},
						"aggs": map[string]interface{}{
							"total_cancelado": map[string]interface{}{
								"sum": map[string]interface{}{
									"field": "cancelamentos.valor_cancelado",
								},
							},
						},
					},
					"soma_negociacoes": map[string]interface{}{
						"nested": map[string]interface{}{
							"path": "negociacoes",
						},
						"aggs": map[string]interface{}{
							"total_negociado": map[string]interface{}{
								"sum": map[string]interface{}{
									"field": "negociacoes.valor_negociado",
								},
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	aggs := result["aggregations"].(map[string]interface{})["resultado"].(map[string]interface{})["buckets"].(map[string]interface{})["all"].(map[string]interface{})

	valorOriginal := aggs["soma_valores_originais"].(map[string]interface{})["value"].(float64)
	cancelamentos := aggs["soma_cancelamentos"].(map[string]interface{})["total_cancelado"].(map[string]interface{})["value"].(float64)
	negociacoes := aggs["soma_negociacoes"].(map[string]interface{})["total_negociado"].(map[string]interface{})["value"].(float64)

	saldoTotal := valorOriginal - cancelamentos - negociacoes

	return map[string]interface{}{
		"codigo_cliente": codigoCliente,
		"periodo": map[string]interface{}{
			"inicio": dataInicio,
			"fim":    dataFim,
		},
		"total_recebiveis": int(hits),
		"saldo_total":      saldoTotal,
		"saldo_formatado":  formatarMoeda(saldoTotal),
	}, nil
}

// Resolver para buscar recebíveis por cliente e data de vencimento
func getReceivablesByCustomerAndDueDateResolver(params graphql.ResolveParams) (interface{}, error) {
	codigoCliente, _ := params.Args["codigo_cliente"].(string)
	dataInicio, _ := params.Args["data_inicio"].(string)
	dataFim, _ := params.Args["data_fim"].(string)
	from, _ := params.Args["from"].(int)
	size, _ := params.Args["size"].(int)

	if size == 0 {
		size = 10
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"codigo_cliente": codigoCliente,
						},
					},
					{
						"range": map[string]interface{}{
							"data_vencimento": map[string]interface{}{
								"gte": dataInicio,
								"lte": dataFim,
							},
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"data_vencimento": map[string]interface{}{
					"order": "asc",
				},
			},
		},
		"from": from,
		"size": size,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	receivables := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		source["id"] = hitMap["_id"]
		receivables[i] = source
	}

	return map[string]interface{}{
		"total":       total,
		"receivables": receivables,
	}, nil
}

// Resolver para contar recebíveis de um cliente
func countReceivablesByCustomerResolver(params graphql.ResolveParams) (interface{}, error) {
	codigoCliente, _ := params.Args["codigo_cliente"].(string)

	if codigoCliente == "" {
		return nil, fmt.Errorf("codigo_cliente é obrigatório")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"codigo_cliente": codigoCliente,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Count(
		esClient.client.Count.WithContext(context.Background()),
		esClient.client.Count.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Count.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	count := int(result["count"].(float64))
	return map[string]interface{}{
		"count": count,
	}, nil
}

// Resolver para contar total de documentos no índice
func getIndexCountResolver(params graphql.ResolveParams) (interface{}, error) {
	res, err := esClient.client.Count(
		esClient.client.Count.WithContext(context.Background()),
		esClient.client.Count.WithIndex("ciclo_vida_recebivel"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	count := int(result["count"].(float64))
	return map[string]interface{}{
		"count": count,
	}, nil
}

// Resolver para contar recebíveis agrupados por cliente
func countReceivablesGroupByCustomerResolver(params graphql.ResolveParams) (interface{}, error) {
	dataInicio, _ := params.Args["data_inicio"].(string)
	dataFim, _ := params.Args["data_fim"].(string)

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"data_vencimento": map[string]interface{}{
					"gte": dataInicio,
					"lte": dataFim,
				},
			},
		},
		"aggs": map[string]interface{}{
			"total_por_cliente": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "codigo_cliente",
					"size":  100,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	aggs := result["aggregations"].(map[string]interface{})["total_por_cliente"].(map[string]interface{})
	buckets := aggs["buckets"].([]interface{})

	stats := make([]map[string]interface{}, len(buckets))
	for i, bucket := range buckets {
		bucketMap := bucket.(map[string]interface{})
		stats[i] = map[string]interface{}{
			"codigo_cliente":   bucketMap["key"].(string),
			"total_recebiveis": int(bucketMap["doc_count"].(float64)),
		}
	}

	return stats, nil
}

// Resolver para buscar cliente com mais registros
func getTopCustomerResolver(params graphql.ResolveParams) (interface{}, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"aggs": map[string]interface{}{
			"clientes": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "codigo_cliente",
					"size":  1,
					"order": map[string]interface{}{
						"_count": "desc",
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Verificar se há agregações
	if aggs, ok := result["aggregations"].(map[string]interface{}); ok {
		if clientes, ok := aggs["clientes"].(map[string]interface{}); ok {
			if buckets, ok := clientes["buckets"].([]interface{}); ok && len(buckets) > 0 {
				bucket := buckets[0].(map[string]interface{})
				return map[string]interface{}{
					"codigo_cliente":   bucket["key"].(string),
					"total_recebiveis": int(bucket["doc_count"].(float64)),
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("nenhum cliente encontrado")
}

// Resolver para buscar recebíveis com saldo disponível mínimo
func getReceivablesByBalanceAvailableResolver(params graphql.ResolveParams) (interface{}, error) {
	codigoCliente, _ := params.Args["codigo_cliente"].(string)
	minBalance, _ := params.Args["min_balance"].(float64)
	from, _ := params.Args["from"].(int)
	size, _ := params.Args["size"].(int)

	if size == 0 {
		size = 50
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"codigo_cliente": codigoCliente,
						},
					},
					{
						"range": map[string]interface{}{
							"valor_original": map[string]interface{}{
								"gte": minBalance,
							},
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"valor_original": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"from":    from,
		"size":    size,
		"_source": []string{"id_recebivel", "codigo_cliente", "valor_original", "data_vencimento", "cancelamentos", "negociacoes"},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	receivables := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		source["id"] = hitMap["_id"]
		receivables[i] = source
	}

	return map[string]interface{}{
		"total":       total,
		"receivables": receivables,
	}, nil
}

// Resolver para buscar saldo de um recebível específico
func getReceivableBalanceByIdResolver(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id é obrigatório")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"_id": id,
			},
		},
		"size": 1,
		"script_fields": map[string]interface{}{
			"saldo_calculado": map[string]interface{}{
				"script": map[string]interface{}{
					"lang":   "painless",
					"source": "double valor = params._source.valor_original; double cancelado = 0.0; double negociado = 0.0; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { cancelado += c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { negociado += n.valor_negociado; } } return Math.round((valor - cancelado - negociado) * 100.0) / 100.0;",
				},
			},
		},
		"_source": []string{"id_recebivel", "codigo_cliente", "valor_original", "data_vencimento", "cancelamentos", "negociacoes"},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.client.Search(
		esClient.client.Search.WithContext(context.Background()),
		esClient.client.Search.WithIndex("ciclo_vida_recebivel"),
		esClient.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return nil, fmt.Errorf("recebível não encontrado")
	}

	hit := hits[0].(map[string]interface{})
	source := hit["_source"].(map[string]interface{})
	source["id"] = hit["_id"]

	if fields, ok := hit["fields"].(map[string]interface{}); ok {
		if saldoArray, ok := fields["saldo_calculado"].([]interface{}); ok && len(saldoArray) > 0 {
			source["saldo_disponivel"] = saldoArray[0].(float64)
		}
	}

	return source, nil
}
