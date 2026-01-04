# 🔷 Guia GraphQL - Data Aggregator

## 📚 Índice
1. [Interface GraphiQL](#interface-graphiql)
2. [Documentação Automática via Introspection](#documentação-automática)
3. [Como Usar](#como-usar)
4. [Exemplos Práticos](#exemplos-práticos)

---

## 🎨 Interface GraphiQL

### O que é GraphiQL?

GraphiQL é uma IDE interativa no navegador para explorar e testar APIs GraphQL. Ela já está **habilitada e funcionando** no seu servidor!

### Como Acessar

1. **Inicie o servidor:**
   ```powershell
   go run .
   ```

2. **Abra o navegador e acesse:**
   ```
   http://localhost:8080/graphql
   ```

### Recursos da Interface GraphiQL

#### ✨ Funcionalidades Principais:

1. **Editor de Queries** (Painel Esquerdo)
   - Syntax highlighting
   - Auto-complete com `Ctrl + Space`
   - Formatação automática com `Ctrl + Shift + P`
   - Validação em tempo real

2. **Painel de Resultados** (Painel Direito)
   - Visualização JSON formatada
   - Mensagens de erro detalhadas
   - Tempo de execução da query

3. **Docs Explorer** (Botão "Docs" no canto superior direito)
   - Documentação completa do schema
   - Navegação por tipos e campos
   - Descrições de queries e argumentos

4. **Query Variables** (Painel inferior)
   - Definir variáveis em JSON
   - Reutilizar queries parametrizadas

5. **Query History** (Histórico de queries executadas)
   - Acessar queries anteriores
   - Reutilizar consultas

---

## 📖 Documentação Automática via Introspection

### O que é Introspection?

GraphQL fornece automaticamente documentação completa do seu schema através de **introspection**. Você não precisa escrever documentação manualmente!

### Como Acessar a Documentação

#### Método 1: Através do GraphiQL (Mais Fácil)

1. Abra `http://localhost:8080/graphql` no navegador
2. Clique no botão **"Docs"** no canto superior direito
3. Navegue pela documentação completa:
   - Lista de todas as queries disponíveis
   - Argumentos obrigatórios e opcionais
   - Tipos de retorno
   - Descrições de cada campo

#### Método 2: Query de Introspection Direta

Você pode consultar o schema programaticamente:

```graphql
# Listar todas as queries disponíveis
query IntrospectionQuery {
  __schema {
    queryType {
      name
      fields {
        name
        description
        args {
          name
          type {
            name
            kind
          }
        }
        type {
          name
          kind
        }
      }
    }
  }
}
```

#### Método 3: Descobrir Tipos Disponíveis

```graphql
query GetAllTypes {
  __schema {
    types {
      name
      kind
      description
      fields {
        name
        type {
          name
          kind
        }
      }
    }
  }
}
```

#### Método 4: Descobrir um Tipo Específico

```graphql
query GetReceivableType {
  __type(name: "Receivable") {
    name
    kind
    description
    fields {
      name
      description
      type {
        name
        kind
      }
    }
  }
}
```

---

## 🚀 Como Usar

### Passo 1: Iniciar o Servidor

```powershell
# No diretório do projeto
cd c:\Users\rafae\OneDrive\Documentos\Workspace\golang\src\data-aggregator

# Iniciar o servidor
go run .
```

Você verá:
```
✅ Conectado ao Elasticsearch com sucesso!
🚀 Servidor HTTP iniciado em http://localhost:8080
📝 Endpoint de query: POST http://localhost:8080/query
💚 Health check: GET http://localhost:8080/health
🔷 GraphQL endpoint: POST http://localhost:8080/graphql
🎨 GraphiQL playground: http://localhost:8080/graphql

✅ Servidor pronto para receber requisições!
```

### Passo 2: Abrir o GraphiQL

Abra seu navegador em: **http://localhost:8080/graphql**

### Passo 3: Explorar a Documentação

1. Clique em **"Docs"** (canto superior direito)
2. Você verá:
   - **Query** - Clique para ver todas as queries disponíveis
   - Cada query mostra:
     - Nome e descrição
     - Argumentos (obrigatórios em negrito)
     - Tipo de retorno

### Passo 4: Escrever Sua Primeira Query

Cole no editor do GraphiQL:

```graphql
query {
  getIndexCount {
    count
  }
}
```

Clique no botão **▶ Execute** (ou pressione `Ctrl + Enter`)

---

## 💡 Exemplos Práticos

### Exemplo 1: Query Simples - Total de Documentos

```graphql
query {
  getIndexCount {
    count
  }
}
```

**Use o Auto-complete:**
1. Digite `{`
2. Pressione `Ctrl + Space`
3. Veja todas as queries disponíveis
4. Selecione `getIndexCount`
5. Dentro de `getIndexCount`, pressione `Ctrl + Space` novamente
6. Veja os campos disponíveis: `count`

### Exemplo 2: Query com Argumentos

```graphql
query {
  countReceivablesByCustomer(codigo_cliente: "CLI-10008") {
    count
  }
}
```

**Dica:** Ao digitar `(` após o nome da query, o GraphiQL mostrará os argumentos disponíveis!

### Exemplo 3: Query com Múltiplos Campos

```graphql
query {
  getReceivableById(id: "REC-00001") {
    id
    id_recebivel
    codigo_cliente
    valor_original
    data_vencimento
    cancelamentos {
      data_cancelamento
      valor_cancelado
      motivo
    }
    negociacoes {
      data_negociacao
      valor_negociado
      tipo_negociacao
    }
  }
}
```

### Exemplo 4: Múltiplas Queries em Uma Requisição

```graphql
query Dashboard {
  totalDocs: getIndexCount {
    count
  }
  
  topCustomer: getTopCustomer {
    codigo_cliente
    total_recebiveis
  }
  
  recentReceivables: getAllReceivables(size: 5) {
    total
    receivables {
      id_recebivel
      codigo_cliente
      valor_original
    }
  }
}
```

### Exemplo 5: Query com Variáveis

**No editor de query:**
```graphql
query GetCustomerData($cliente: String!, $inicio: String!, $fim: String!) {
  getCustomerBalance(
    codigo_cliente: $cliente
    data_inicio: $inicio
    data_fim: $fim
  ) {
    codigo_cliente
    total_recebiveis
    saldo_total
    saldo_formatado
    periodo {
      inicio
      fim
    }
  }
}
```

**No painel "Query Variables" (abaixo do editor):**
```json
{
  "cliente": "123456",
  "inicio": "2025-01-01",
  "fim": "2026-12-31"
}
```

### Exemplo 6: Fragmentos para Reutilização

```graphql
fragment ReceivableBasicInfo on Receivable {
  id
  id_recebivel
  codigo_cliente
  valor_original
  data_vencimento
}

query {
  receivable1: getReceivableById(id: "REC-00001") {
    ...ReceivableBasicInfo
  }
  
  receivable2: getReceivableById(id: "REC-00002") {
    ...ReceivableBasicInfo
  }
}
```

---

## 🔍 Explorando a Documentação no GraphiQL

### Visualizando Todas as Queries

1. Abra o GraphiQL
2. Clique em **"Docs"**
3. Clique em **"Query"**
4. Você verá todas as 10 queries:

#### Lista Completa de Queries:

1. **getAllReceivables(size: Int = 10): SearchResult**
   - Buscar todos os recebíveis com limite
   - Argumento opcional: `size` (padrão: 10)

2. **getReceivableById(id: String!): Receivable**
   - Buscar recebível por ID
   - Argumento obrigatório: `id`

3. **getCustomerBalance(codigo_cliente: String!, data_inicio: String!, data_fim: String!): Balance**
   - Buscar saldo de um cliente por período
   - Argumentos obrigatórios: `codigo_cliente`, `data_inicio`, `data_fim`

4. **getReceivablesByCustomerAndDueDate(...): SearchResult**
   - Buscar recebíveis por cliente e data de vencimento
   - Com paginação (`from`, `size`)

5. **countReceivablesByCustomer(codigo_cliente: String!): CountResult**
   - Contar recebíveis de um cliente específico

6. **getIndexCount(): CountResult**
   - Contar total de documentos no índice

7. **countReceivablesGroupByCustomer(data_inicio: String, data_fim: String): [CustomerStats]**
   - Contar recebíveis agrupados por cliente

8. **getTopCustomer(): CustomerStats**
   - Buscar cliente com mais registros

9. **getReceivablesByBalanceAvailable(...): SearchResult**
   - Buscar recebíveis com saldo disponível mínimo

10. **getReceivableBalanceById(id: String!): Receivable**
    - Buscar saldo de um recebível específico por ID

### Visualizando Tipos

Clique em qualquer tipo (ex: `Receivable`, `Balance`) para ver:
- Todos os campos disponíveis
- Tipo de cada campo
- Se é obrigatório (!) ou opcional

---

## 🎯 Dicas e Truques

### Atalhos do Teclado

| Atalho | Ação |
|--------|------|
| `Ctrl + Enter` | Executar query |
| `Ctrl + Space` | Auto-complete |
| `Ctrl + Shift + P` | Formatar query |
| `Ctrl + /` | Comentar/descomentar linha |

### Formatar Query Automaticamente

Cole uma query sem formatação:
```graphql
query{getAllReceivables(size:5){total receivables{id codigo_cliente}}}
```

Pressione `Ctrl + Shift + P` e ela será formatada automaticamente:
```graphql
query {
  getAllReceivables(size: 5) {
    total
    receivables {
      id
      codigo_cliente
    }
  }
}
```

### Validação em Tempo Real

O GraphiQL valida sua query enquanto você digita:
- ✅ Verde: Query válida
- ❌ Vermelho: Erro de sintaxe ou campo inexistente
- Passe o mouse sobre o erro para ver detalhes

### Descobrir Campos Aninhados

Ao digitar `{` após um campo, pressione `Ctrl + Space` para ver os subcampos disponíveis:

```graphql
query {
  getReceivableById(id: "REC-00001") {
    cancelamentos {
      # Pressione Ctrl + Space aqui para ver:
      # - data_cancelamento
      # - valor_cancelado
      # - motivo
    }
  }
}
```

---

## 🧪 Testando a Introspection

### Teste 1: Listar Todas as Queries

```graphql
query {
  __schema {
    queryType {
      fields {
        name
        description
      }
    }
  }
}
```

### Teste 2: Ver Detalhes de Uma Query

```graphql
query {
  __type(name: "Query") {
    fields {
      name
      description
      args {
        name
        type {
          name
          kind
        }
      }
    }
  }
}
```

### Teste 3: Explorar Tipo Receivable

```graphql
query {
  __type(name: "Receivable") {
    name
    fields {
      name
      type {
        name
        kind
        ofType {
          name
          kind
        }
      }
    }
  }
}
```

---

## 📱 Testando via HTTP (sem GraphiQL)

Se preferir usar curl, Postman ou VS Code REST Client:

### Exemplo com curl:

```powershell
curl -X POST http://localhost:8080/graphql `
  -H "Content-Type: application/json" `
  -d '{\"query\": \"{ getIndexCount { count } }\"}'
```

### Exemplo com arquivo .http:

Já existe o arquivo `requests/graphql_queries.http` com exemplos prontos!

---

## 🎓 Próximos Passos

Agora que você conhece:
- ✅ Interface GraphiQL
- ✅ Documentação automática via introspection

Você pode:

1. **Explorar** todas as queries no GraphiQL
2. **Testar** diferentes combinações de campos
3. **Criar** queries customizadas para suas necessidades
4. **Usar** variáveis para queries dinâmicas
5. **Combinar** múltiplas queries em uma requisição
6. **Integrar** o GraphQL em suas aplicações frontend

---

## 🆘 Solução de Problemas

### GraphiQL não carrega?

1. Verifique se o servidor está rodando
2. Acesse `http://localhost:8080/health` para testar
3. Verifique se não há firewall bloqueando a porta 8080

### Auto-complete não funciona?

1. Certifique-se de que a query está sintaticamente correta até o ponto do cursor
2. Use `Ctrl + Space` após digitar `{` ou após o nome de um campo

### Query retorna erro?

1. Verifique a aba "Docs" para ver os argumentos obrigatórios
2. Use o auto-complete para evitar erros de digitação
3. Veja detalhes do erro no painel de resultados

---

## 📚 Recursos Adicionais

- [Documentação GraphQL](https://graphql.org/learn/)
- [GraphiQL Documentation](https://github.com/graphql/graphiql)
- [GraphQL Introspection](https://graphql.org/learn/introspection/)

---

**Desenvolvido para Data Aggregator - Elasticsearch + GraphQL** 🚀
