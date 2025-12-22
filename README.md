# Data Aggregator - API REST para Elasticsearch

API REST em Go para gerenciar documentos no Elasticsearch com operações CRUD completas através de um endpoint genérico.

## 🚀 Funcionalidades

- ✅ **API REST** com endpoint único `/query` para todas as operações
- ✅ **Criar índices** com mapeamentos customizados
- ✅ **Inserir documentos** com ID específico ou auto-gerado
- ✅ **Atualizar documentos** parcialmente
- ✅ **Buscar documentos** por ID
- ✅ **Buscar documentos** com queries complexas (match, range, bool, etc)
- ✅ **Deletar documentos**
- ✅ **Health check** endpoint

## 📋 Pré-requisitos

- Go 1.23.4 ou superior
- Docker e Docker Compose

## 🔧 Configuração

### 1. Inicie o Docker Desktop

Certifique-se que o Docker Desktop está rodando no Windows.

### 2. Crie a rede Docker

```powershell
docker network create data-aggregator-network
```

### 3. Inicie o Elasticsearch

```powershell
docker-compose up -d
```

### 4. Verifique se o Elasticsearch está rodando

```powershell
curl http://localhost:9200
```

Você deve ver uma resposta JSON com informações do cluster.

## ▶️ Executar a Aplicação

```powershell
go run main.go
```

A API estará disponível em `http://localhost:8080`

## 📡 API Endpoints

### Health Check

Verifica se o serviço está rodando.

```powershell
curl http://localhost:8080/health
```

**Resposta:**
```json
{
  "status": "ok",
  "service": "data-aggregator",
  "time": "2025-12-21T18:00:00Z"
}
```

### Query Endpoint (Genérico)

Endpoint único para todas as operações do Elasticsearch.

**URL:** `POST http://localhost:8080/query`

**Formato da Requisição:**
```json
{
  "operation": "nome_da_operacao",
  "index": "nome_do_indice",
  "document_id": "id_do_documento",
  "body": { ... }
}
```

**Operações disponíveis:**
- `create_index` - Criar índice
- `index` - Inserir documento
- `update` - Atualizar documento
- `get` - Buscar documento por ID
- `search` - Buscar documentos com query
- `delete` - Deletar documento

---

## 📝 Exemplos de Uso

### 1. Criar Índice

Cria um novo índice com mapeamento de campos.

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "create_index",
    "index": "products",
    "body": {
      "mappings": {
        "properties": {
          "name": { "type": "text" },
          "description": { "type": "text" },
          "price": { "type": "float" },
          "quantity": { "type": "integer" },
          "category": { "type": "keyword" },
          "timestamp": { "type": "date" }
        }
      }
    }
  }'
```

**Resposta de Sucesso:**
```json
{
  "success": true,
  "message": "Índice 'products' criado com sucesso"
}
```

---

### 2. Inserir Documento

Insere um novo documento no índice.

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "index",
    "index": "products",
    "document_id": "prod-1",
    "body": {
      "name": "Laptop Dell XPS 15",
      "description": "Notebook de alta performance",
      "price": 8999.99,
      "quantity": 10,
      "category": "electronics",
      "timestamp": "2025-12-21T18:00:00Z"
    }
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Documento 'prod-1' inserido com sucesso"
}
```

**Inserir múltiplos documentos:**
---

## 🎯 Exemplos de Casos de Uso

### Criar um E-commerce de Produtos

```powershell
# 1. Criar índice
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "create_index",
  "index": "products",
  "body": {
    "mappings": {
      "properties": {
        "name": {"type": "text"},
        "price": {"type": "float"},
        "category": {"type": "keyword"},
        "in_stock": {"type": "boolean"}
      }
    }
  }
}'

# 2. Adicionar produtos
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "index",
  "index": "products",
  "document_id": "prod-1",
  "body": {"name": "iPhone 15", "price": 6999.00, "category": "smartphones", "in_stock": true}
}'

# 3. Buscar produtos em estoque
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "products",
  "body": {
    "query": {"match": {"in_stock": true}}
  }
}'

# 4. Atualizar preço
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "update",
  "index": "products",
  "document_id": "prod-1",
  "body": {"price": 6499.00}
}'
```

### Ciclo de Vida de Recebível

Sistema de gerenciamento de recebíveis com cancelamentos e negociações.

```powershell
# 1. Criar índice de recebíveis
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "create_index",
  "index": "ciclo_vida_recebivel",
  "body": {
    "mappings": {
      "properties": {
        "id_recebivel": {"type": "keyword"},
        "codigo_cliente": {"type": "keyword"},
        "codigo_produto": {"type": "integer"},
        "codigo_produto_parceiro": {"type": "integer"},
        "modalidade": {"type": "integer"},
        "valor_original": {"type": "float"},
        "data_vencimento": {"type": "date"},
        "cancelamentos": {
          "type": "nested",
          "properties": {
            "id_cancelamento": {"type": "keyword"},
            "data_cancelamento": {"type": "date"},
            "valor_cancelado": {"type": "float"},
            "motivo": {"type": "text"}
          }
        },
        "negociacoes": {
          "type": "nested",
          "properties": {
            "id_negociacao": {"type": "keyword"},
            "data_negociacao": {"type": "date"},
            "valor_negociado": {"type": "float"}
          }
        }
      }
    }
  }
}'

# 2. Inserir recebível com cancelamentos e negociações
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "index",
  "index": "ciclo_vida_recebivel",
  "document_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
  "body": {
    "id_recebivel": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
    "codigo_cliente": "CLI-12345",
    "codigo_produto": 123,
    "codigo_produto_parceiro": 25,
    "modalidade": 1,
    "valor_original": 10000.00,
    "data_vencimento": "2025-12-31",
    "cancelamentos": [
      {
        "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7c",
        "data_cancelamento": "2025-12-21",
        "valor_cancelado": 1000.00,
        "motivo": "Cliente solicitou cancelamento parcial."
      },
      {
        "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7d",
        "data_cancelamento": "2025-12-12",
        "valor_cancelado": 500.00,
        "motivo": "Cliente solicitou cancelamento parcial."
      }
    ],
    "negociacoes": [
      {
        "id_negociacao": "d4e5f6a7-b8c9-4d0e-9f1a-2b3c4d5e6f7g",
        "data_negociacao": "2025-12-19",
        "valor_negociado": 500.00
      }
    ]
  }
}'

# 3. Buscar recebível por ID
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "get",
  "index": "ciclo_vida_recebivel",
  "document_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c"
}'

# 4. Buscar recebíveis com cancelamentos
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "nested": {
        "path": "cancelamentos",
        "query": {
          "range": {
            "cancelamentos.valor_cancelado": {"gte": 500}
          }
        }
      }
    }
  }
}'

# 5. Buscar recebíveis por modalidade
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "term": {"modalidade": 1}
    }
  }
}'

# 6. Buscar recebíveis por cliente e período de vencimento
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "bool": {
        "must": [
          {
            "term": {
              "codigo_cliente": "CLI-12345"
            }
          },
          {
            "range": {
              "data_vencimento": {
                "gte": "2025-12-01",
                "lte": "2025-12-31"
              }
            }
          }
        ]
      }
    },
    "sort": [
      {
        "data_vencimento": {
          "order": "asc"
        }
      }
    ]
  }
}'

# 7. Buscar recebíveis por cliente e período com saldo individual e total
# Retorna o saldo de cada recebível + soma total dos saldos
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "bool": {
        "must": [
          {
            "term": {
              "codigo_cliente": "CLI-12345"
            }
          },
          {
            "range": {
              "data_vencimento": {
                "gte": "2025-12-01",
                "lte": "2025-12-31"
              }
            }
          }
        ]
      }
    },
    "script_fields": {
      "saldo": {
        "script": {
          "lang": "painless",
          "source": "double saldo = doc[\"valor_original\"].value; if (params._source.containsKey(\"cancelamentos\") && params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.containsKey(\"negociacoes\") && params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } return saldo;"
        }
      }
    },
    "_source": ["id_recebivel", "data_vencimento"],
    "aggs": {
      "saldo_total": {
        "scripted_metric": {
          "init_script": "state.saldo_total = 0.0",
          "map_script": "double saldo = doc[\"valor_original\"].value; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } state.saldo_total += saldo;",
          "combine_script": "return state.saldo_total",
          "reduce_script": "double total = 0; for (s in states) { total += s; } return total"
        }
      }
    },
    "sort": [
      {
        "data_vencimento": {
          "order": "asc"
        }
      }
    ]
  }
}'
```

**Explicação:**
Esta consulta combina:
- **Filtros**: `codigo_cliente` + período de `data_vencimento`
- **script_fields**: Calcula o saldo individual de cada recebível
- **aggregations**: Calcula a soma total de todos os saldos encontrados

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Encontrados 1 documentos",
  "data": {
    "hits": [
      {
        "_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
        "_source": {
          "id_recebivel": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
          "data_vencimento": "2025-12-31"
        },
        "fields": {
          "saldo": [8000.00]
        }
      }
    ],
    "aggregations": {
      "saldo_total": {
        "value": 8000.00
      }
    }
  }
}
```

Neste exemplo:
- **Saldo individual do recebível**: R$ 8.000,00
- **Saldo total de todos os recebíveis do cliente CLI-12345 no período**: R$ 8.000,00

# 8. Adicionar novo cancelamento (atualizar documento)
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "update",
  "index": "ciclo_vida_recebivel",
  "document_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
  "body": {
    "cancelamentos": [
      {
        "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7c",
        "data_cancelamento": "2025-12-21",
        "valor_cancelado": 1000.00,
        "motivo": "Cliente solicitou cancelamento parcial."
      },
      {
        "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7d",
        "data_cancelamento": "2025-12-12",
        "valor_cancelado": 500.00,
        "motivo": "Cliente solicitou cancelamento parcial."
      },
      {
        "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7e",
        "data_cancelamento": "2025-12-22",
        "valor_cancelado": 300.00,
        "motivo": "Ajuste de valor."
      }
    ]
  }
}'

# 9. Buscar recebíveis com cálculo de saldo
# Saldo = valor_original - SUM(valor_cancelado) - SUM(valor_negociado)
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "match_all": {}
    },
    "script_fields": {
      "saldo": {
        "script": {
          "lang": "painless",
          "source": "double saldo = doc[\"valor_original\"].value; if (params._source.containsKey(\"cancelamentos\") && params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.containsKey(\"negociacoes\") && params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } return saldo;"
        }
      }
    },
    "_source": ["id_recebivel", "valor_original", "cancelamentos", "negociacoes"]
  }
}'
```

**Explicação do Cálculo de Saldo:**

O script Painless realiza o seguinte cálculo:
```
Saldo = valor_original - Σ(valor_cancelado) - Σ(valor_negociado)
```

Para o exemplo do recebível `a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c`:
- Valor Original: R$ 10.000,00
- Cancelamentos: R$ 1.000,00 + R$ 500,00 = R$ 1.500,00
- Negociações: R$ 500,00
- **Saldo Final: R$ 10.000,00 - R$ 1.500,00 - R$ 500,00 = R$ 8.000,00**

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Encontrados 1 documentos",
  "data": {
    "hits": [
      {
        "_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
        "_source": {
          "id_recebivel": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c"
        },
        "fields": {
          "saldo": [8000.00]
        }
      }
    ]
  }
}
```

# 10. Buscar recebível específico com cálculo de saldo (por ID)
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "query": {
      "term": {
        "id_recebivel": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c"
      }
    },
    "script_fields": {
      "saldo": {
        "script": {
          "lang": "painless",
          "source": "double saldo = doc[\"valor_original\"].value; if (params._source.containsKey(\"cancelamentos\") && params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.containsKey(\"negociacoes\") && params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } return saldo;"
        }
      },
      "total_cancelado": {
        "script": {
          "lang": "painless",
          "source": "double total = 0; if (params._source.containsKey(\"cancelamentos\") && params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { total += c.valor_cancelado; } } return total;"
        }
      },
      "total_negociado": {
        "script": {
          "lang": "painless",
          "source": "double total = 0; if (params._source.containsKey(\"negociacoes\") && params._source.negociacoes != null) { for (def n : params._source.negociacoes) { total += n.valor_negociado; } } return total;"
        }
      }
    },
    "_source": ["id_recebivel", "codigo_produto", "modalidade", "valor_original", "data_vencimento", "cancelamentos", "negociacoes"]
  }
}'
```

**Resposta esperada para busca por ID:**
```json
{
  "success": true,
  "message": "Encontrados 1 documentos",
  "data": {
    "hits": [
      {
        "_id": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
        "_source": {
          "id_recebivel": "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c",
          "codigo_produto": 123,
          "modalidade": 1,
          "valor_original": 10000.00,
          "data_vencimento": "2025-12-31",
          "cancelamentos": [
            {
              "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7c",
              "data_cancelamento": "2025-12-21",
              "valor_cancelado": 1000.00,
              "motivo": "Cliente solicitou cancelamento parcial."
            },
            {
              "id_cancelamento": "c3d9f8e4-2b6a-4f7b-9d3e-6f8b0f2b4d7d",
              "data_cancelamento": "2025-12-12",
              "valor_cancelado": 500.00,
              "motivo": "Cliente solicitou cancelamento parcial."
            }
          ],
          "negociacoes": [
            {
              "id_negociacao": "d4e5f6a7-b8c9-4d0e-9f1a-2b3c4d5e6f7g",
              "data_negociacao": "2025-12-19",
              "valor_negociado": 500.00
            }
          ]
        },
        "fields": {
          "saldo": [8000.00],
          "total_cancelado": [1500.00],
          "total_negociado": [500.00]
        }
      }
    ]
  }
}
```

Esta consulta retorna:
- **saldo**: R$ 8.000,00 (valor líquido após cancelamentos e negociações)
- **total_cancelado**: R$ 1.500,00 (soma de todos os cancelamentos)
- **total_negociado**: R$ 500,00 (soma de todas as negociações)

**Versão alternativa com aggregation (para múltiplos recebíveis):**

```powershell
# Buscar com agregação de saldo médio por modalidade
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{
  "operation": "search",
  "index": "ciclo_vida_recebivel",
  "body": {
    "size": 0,
    "aggs": {
      "por_modalidade": {
        "terms": {
          "field": "modalidade"
        },
        "aggs": {
          "saldo_calculado": {
            "scripted_metric": {
              "init_script": "state.saldos = []",
              "map_script": "double saldo = doc[\"valor_original\"].value; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } state.saldos.add(saldo);",
              "combine_script": "return state.saldos",
              "reduce_script": "double total = 0; for (s in states) { for (saldo in s) { total += saldo; } } return total"
            }
          }
        }
      }
    }
  }
}'
```

---

## 🔍 Referência de Query DSL

### Operadores de Range

- `gte` - Greater than or equal (maior ou igual)
- `gt` - Greater than (maior que)
- `lte` - Less than or equal (menor ou igual)
- `lt` - Less than (menor que)

### Tipos de Query

- **match** - Busca por texto com análise
- **term** - Busca exata (keywords)
- **range** - Busca por intervalo de valores
- **bool** - Combina múltiplas queries
  - `must` - Deve corresponder (AND)
  - `should` - Pode corresponder (OR)
  - `must_not` - Não deve corresponder (NOT)
  - `filter` - Deve corresponder (sem score)
- **wildcard** - Busca com curingas (* e ?)
- **match_all** - Retorna todos os documentos

---

## 🐛 Troubleshooting

### Erro: "connection refused"

Verifique se o Elasticsearch está rodando:
```powershell
docker ps | Select-String elasticsearch
```

Se não estiver, inicie:
```powershell
docker-compose up -d
```

### Erro: "index already exists"

Use uma requisição de update ou delete o índice existente primeiro através do Kibana ou comandos curl diretos ao Elasticsearch.

### Verificar saúde do cluster

```powershell
curl http://localhost:9200/_cluster/health?pretty
```

### Verificar índices existentes

```powershell
curl http://localhost:9200/_cat/indices?v
```

### API não responde

Verifique se a aplicação está rodando:
```powershell
curl http://localhost:8080/health
```

---

## 📚 Recursos

- [Elasticsearch Go Client Documentation](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)
- [Elasticsearch Query DSL](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html)
- [Elasticsearch Mapping](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html)
- [Elasticsearch REST API](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html)

---

## 🎯 Próximos Passos

- [ ] Adicionar autenticação JWT
- [ ] Implementar rate limiting
- [ ] Adicionar suporte a bulk operations
- [ ] Implementar aggregations
- [ ] Adicionar validação de dados
- [ ] Criar testes unitários e de integração
- [ ] Adicionar logging estruturado
- [ ] Implementar circuit breaker
- [ ] Adicionar métricas Prometheus
  "index": "app-logs",
  "body": {
    "query": {"match": {"level": "error"}},
    "sort": [{"timestamp": {"order": "desc"}}],
    "size": 100
  }
}'
```H "Content-Type: application/json" `
  -d '{
    "operation": "get",
    "index": "products",
    "document_id": "prod-1"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Documento encontrado",
  "data": {
    "_id": "prod-1",
    "_index": "products",
    "_source": {
      "name": "Laptop Dell XPS 15",
      "price": 8999.99,
      "quantity": 10,
      "category": "electronics"
    }
  }
}
```

---

### 4. Atualizar Documento

Atualiza campos específicos de um documento existente.

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "update",
    "index": "products",
    "document_id": "prod-1",
    "body": {
      "price": 7999.99,
      "quantity": 8
    }
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Documento 'prod-1' atualizado com sucesso"
}
```

---

### 5. Buscar Documentos com Query

Busca documentos usando queries do Elasticsearch.

#### Buscar por Match (Texto)

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "search",
    "index": "products",
    "body": {
      "query": {
        "match": {
          "category": "electronics"
        }
      },
      "size": 10
    }
  }'
```

#### Buscar por Range (Preço)

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "search",
    "index": "products",
    "body": {
      "query": {
        "range": {
          "price": {
            "gte": 500,
            "lte": 9000
          }
        }
      }
    }
  }'
```

#### Buscar com Ordenação

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "search",
    "index": "products",
    "body": {
      "query": {
        "match_all": {}
      },
      "sort": [
        { "price": { "order": "desc" } }
      ]
    }
  }'
```

#### Buscar com Bool Query (Múltiplas Condições)

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "search",
    "index": "products",
    "body": {
      "query": {
        "bool": {
          "must": [
            { "match": { "category": "electronics" } }
          ],
          "filter": [
            { "range": { "price": { "gte": 500 } } }
          ]
        }
      }
    }
  }'
```

#### Buscar por Texto Parcial (Wildcard)

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "search",
    "index": "products",
    "body": {
      "query": {
        "wildcard": {
          "name": "*laptop*"
        }
      }
    }
  }'
```

**Resposta de Search:**
```json
{
  "success": true,
  "message": "Encontrados 3 documentos",
  "data": {
    "hits": [
      {
        "_id": "prod-1",
        "_score": 1.0,
        "_source": {
          "name": "Laptop Dell XPS 15",
          "price": 7999.99
        }
      }
    ],
    "total": 3
  }
}
```

---

### 6. Deletar Documento

Remove um documento do índice.

```powershell
curl -X POST http://localhost:8080/query `
  -H "Content-Type: application/json" `
  -d '{
    "operation": "delete",
    "index": "products",
    "document_id": "prod-3"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Documento 'prod-3' deletado com sucesso"
}
```

## 🔍 Queries Avançadas

### Busca por Range (Intervalo)

```go
query := map[string]interface{}{
    "query": map[string]interface{}{
        "range": map[string]interface{}{
            "price": map[string]interface{}{
                "gte": 100,  // maior ou igual
                "lte": 500,  // menor ou igual
            },
        },
    },
}
```

### Busca Bool (Múltiplas condições)

```go
query := map[string]interface{}{
    "query": map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {"match": map[string]interface{}{"category": "electronics"}},
                {"range": map[string]interface{}{
                    "price": map[string]interface{}{"gte": 100},
                }},
            },
        },
    },
}
```

### Busca com Ordenação

```go
query := map[string]interface{}{
    "query": map[string]interface{}{
        "match_all": map[string]interface{}{},
    },
    "sort": []map[string]interface{}{
        {"price": map[string]string{"order": "desc"}},
    },
}
```

## 🐛 Troubleshooting

### Erro: "connection refused"

Verifique se o Elasticsearch está rodando:
```powershell
docker ps | Select-String elasticsearch
```

Se não estiver, inicie:
```powershell
docker-compose up -d
```

### Erro: "index already exists"

Se o índice já existir, a aplicação irá continuar normalmente com um aviso.

### Verificar saúde do cluster

```powershell
curl http://localhost:9200/_cluster/health?pretty
```

---

## 🔗 Queries Diretas no Elasticsearch

Todos os exemplos acima podem ser executados diretamente no Elasticsearch, sem usar a API REST. Aqui estão as conversões:

### 1. Criar Índice (Produtos)

```powershell
curl -X PUT "http://localhost:9200/products" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "name": {"type": "text"},
      "description": {"type": "text"},
      "price": {"type": "float"},
      "quantity": {"type": "integer"},
      "category": {"type": "keyword"},
      "timestamp": {"type": "date"}
    }
  }
}'
```

### 2. Criar Índice (Recebíveis)

```powershell
curl -X PUT "http://localhost:9200/ciclo_vida_recebivel" -H "Content-Type: application/json" -d '{
  "mappings": {
    "properties": {
      "id_recebivel": {"type": "keyword"},
      "codigo_cliente": {"type": "keyword"},
      "codigo_produto": {"type": "integer"},
      "codigo_produto_parceiro": {"type": "integer"},
      "modalidade": {"type": "integer"},
      "valor_original": {"type": "float"},
      "data_vencimento": {"type": "date"},
      "cancelamentos": {
        "type": "nested",
        "properties": {
          "id_cancelamento": {"type": "keyword"},
          "data_cancelamento": {"type": "date"},
          "valor_cancelado": {"type": "float"},
          "motivo": {"type": "text"}
        }
      },
      "negociacoes": {
        "type": "nested",
        "properties": {
          "id_negociacao": {"type": "keyword"},
          "data_negociacao": {"type": "date"},
          "valor_negociado": {"type": "float"}
        }
      }
    }
  }
}'
```

### 3. Inserir Documento

```powershell
curl -X PUT "http://localhost:9200/products/_doc/prod-1" -H "Content-Type: application/json" -d '{
  "name": "Laptop Dell XPS 15",
  "description": "Notebook de alta performance",
  "price": 8999.99,
  "quantity": 10,
  "category": "electronics",
  "timestamp": "2025-12-21T18:00:00Z"
}'
```

### 4. Buscar Documento por ID

```powershell
curl -X GET "http://localhost:9200/products/_doc/prod-1"
```

### 5. Atualizar Documento

```powershell
curl -X POST "http://localhost:9200/products/_update/prod-1" -H "Content-Type: application/json" -d '{
  "doc": {
    "price": 7999.99,
    "quantity": 8
  }
}'
```

### 6. Buscar com Query (Match)

```powershell
curl -X GET "http://localhost:9200/products/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "match": {
      "category": "electronics"
    }
  },
  "size": 10
}'
```

### 7. Buscar com Range

```powershell
curl -X GET "http://localhost:9200/products/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "range": {
      "price": {
        "gte": 500,
        "lte": 9000
      }
    }
  }
}'
```

### 8. Buscar com Bool Query (Múltiplas Condições)

```powershell
curl -X GET "http://localhost:9200/products/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "bool": {
      "must": [
        {"match": {"category": "electronics"}}
      ],
      "filter": [
        {"range": {"price": {"gte": 500}}}
      ]
    }
  }
}'
```

### 9. Buscar Recebíveis com Cancelamentos (Nested Query)

```powershell
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "nested": {
      "path": "cancelamentos",
      "query": {
        "range": {
          "cancelamentos.valor_cancelado": {"gte": 500}
        }
      }
    }
  }
}'
```

### 10. Buscar Recebíveis por Cliente e Período

```powershell
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "bool": {
      "must": [
        {"term": {"codigo_cliente": "CLI-12345"}},
        {"range": {"data_vencimento": {"gte": "2025-12-01", "lte": "2025-12-31"}}}
      ]
    }
  },
  "sort": [{"data_vencimento": {"order": "asc"}}]
}'
```

### 11. Buscar Recebíveis com Saldo Calculado

```powershell
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "bool": {
      "must": [
        {"term": {"codigo_cliente": "CLI-12345"}},
        {"range": {"data_vencimento": {"gte": "2025-12-01", "lte": "2025-12-31"}}}
      ]
    }
  },
  "script_fields": {
    "saldo": {
      "script": {
        "lang": "painless",
        "source": "double saldo = doc[\"valor_original\"].value; if (params._source.containsKey(\"cancelamentos\") && params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.containsKey(\"negociacoes\") && params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } return saldo;"
      }
    }
  },
  "_source": ["id_recebivel", "data_vencimento"],
  "aggs": {
    "saldo_total": {
      "scripted_metric": {
        "init_script": "state.saldo_total = 0.0",
        "map_script": "double saldo = doc[\"valor_original\"].value; if (params._source.cancelamentos != null) { for (def c : params._source.cancelamentos) { saldo -= c.valor_cancelado; } } if (params._source.negociacoes != null) { for (def n : params._source.negociacoes) { saldo -= n.valor_negociado; } } state.saldo_total += saldo;",
        "combine_script": "return state.saldo_total",
        "reduce_script": "double total = 0; for (s in states) { total += s; } return total"
      }
    }
  },
  "sort": [{"data_vencimento": {"order": "asc"}}]
}'
```

### 12. Deletar Documento

```powershell
curl -X DELETE "http://localhost:9200/products/_doc/prod-3"
```

### 13. Contar Documentos

```powershell
# Contar todos os documentos
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_count"

# Contar com filtro
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_count" -H "Content-Type: application/json" -d '{
  "query": {
    "term": {"codigo_cliente": "CLI-10001"}
  }
}'
```

### 14. Verificar Índices Existentes

```powershell
curl -X GET "http://localhost:9200/_cat/indices?v"
```

### 15. Deletar Índice

```powershell
curl -X DELETE "http://localhost:9200/ciclo_vida_recebivel"
```

### 16. Estatísticas do Índice

```powershell
curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_stats"
```

---

## 📚 Recursos

- [Elasticsearch Go Client Documentation](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)
- [Elasticsearch Query DSL](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html)
- [Elasticsearch Mapping](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html)
- [Elasticsearch REST API](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html)

## 🎯 Próximos Passos

- [ ] Adicionar suporte a bulk operations
- [ ] Implementar aggregations
- [ ] Adicionar validação de dados
- [ ] Criar API REST com handlers HTTP
- [ ] Adicionar testes unitários
# aggregator
