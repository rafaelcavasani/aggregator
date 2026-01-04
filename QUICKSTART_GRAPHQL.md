# 🚀 Quick Start - GraphQL + Introspection

## ✅ Os Dois Recursos Já Estão Implementados!

Ótimas notícias! Tanto a **Interface GraphiQL** quanto a **Documentação Automática via Introspection** já estão funcionando no seu projeto. Veja como usar:

---

## 🎨 1. Interface GraphiQL

### O que é?
GraphiQL é uma IDE web interativa para testar e explorar APIs GraphQL. Ela está **ativa e pronta para uso**!

### Como acessar?

**Passo 1:** Inicie o servidor
```powershell
cd c:\Users\rafae\OneDrive\Documentos\Workspace\golang\src\data-aggregator
go run .
```

**Passo 2:** Abra o navegador em:
```
http://localhost:8080/graphql
```

### O que você verá?

```
┌─────────────────────────────────────────────────┐
│  🔷 GraphiQL Interface                          │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐  ┌───────────────────┐       │
│  │  Editor      │  │  Resultados       │       │
│  │  (Esquerda)  │  │  (Direita)        │       │
│  │              │  │                   │       │
│  │  Digite suas │  │  Veja os          │       │
│  │  queries     │  │  resultados       │       │
│  │  aqui        │  │  aqui             │       │
│  └──────────────┘  └───────────────────┘       │
│                                                 │
│  [Docs] [History] [▶ Execute]                  │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Funcionalidades Principais:

#### ✨ Auto-complete
- Digite `{` e pressione `Ctrl + Space`
- Veja todas as queries disponíveis
- Selecione com as setas e Enter

#### 📖 Documentação Integrada
- Clique no botão **"Docs"** (canto superior direito)
- Navegue por todas as queries
- Veja argumentos e tipos de retorno

#### ⚡ Execução Rápida
- Pressione `Ctrl + Enter` para executar
- Ou clique no botão ▶ Execute

#### 🎨 Formatação Automática
- Pressione `Ctrl + Shift + P`
- Sua query será formatada automaticamente

---

## 📚 2. Documentação Automática via Introspection

### O que é?
GraphQL fornece documentação automática através do sistema de **introspection**. Você pode consultar o próprio schema para descobrir o que está disponível!

### Como usar?

#### Método 1: Através do GraphiQL (Recomendado)

1. Abra `http://localhost:8080/graphql`
2. Clique em **"Docs"** no canto superior direito
3. Explore:
   - **Query** → Ver todas as queries
   - **Tipos** → Ver estruturas de dados
   - **Argumentos** → Ver parâmetros obrigatórios

#### Método 2: Queries de Introspection

##### 📋 Listar todas as queries disponíveis:

```graphql
{
  __schema {
    queryType {
      fields {
        name
        description
        args {
          name
          type {
            name
          }
        }
      }
    }
  }
}
```

##### 🔍 Ver detalhes de um tipo específico:

```graphql
{
  __type(name: "Receivable") {
    name
    kind
    fields {
      name
      type {
        name
        kind
      }
    }
  }
}
```

##### 📊 Listar todos os tipos disponíveis:

```graphql
{
  __schema {
    types {
      name
      kind
      description
    }
  }
}
```

---

## 🧪 Teste Agora!

### Teste 1: Query Simples
Cole no GraphiQL e execute (Ctrl + Enter):

```graphql
query {
  getIndexCount {
    count
  }
}
```

**Resultado esperado:**
```json
{
  "data": {
    "getIndexCount": {
      "count": 1000
    }
  }
}
```

### Teste 2: Query com Auto-complete
1. Digite apenas: `{`
2. Pressione `Ctrl + Space`
3. Veja a lista de queries
4. Selecione `getTopCustomer`
5. Digite `{` novamente dentro
6. Pressione `Ctrl + Space`
7. Selecione os campos desejados

### Teste 3: Explorar Documentação
1. Clique em **"Docs"**
2. Clique em **"Query"**
3. Clique em **"getCustomerBalance"**
4. Veja:
   - Descrição: "Buscar saldo de um cliente por período"
   - Argumentos obrigatórios em **negrito**
   - Tipo de retorno: Balance

### Teste 4: Introspection via Script
Execute o script PowerShell:

```powershell
.\test_introspection.ps1
```

Escolha uma opção do menu interativo!

---

## 📖 Recursos Criados

Todos esses arquivos foram criados para você:

### 📄 Documentação
- **GRAPHQL_GUIDE.md** → Guia completo e detalhado
- **graphql_queries_examples.md** → Exemplos de todas as queries
- **QUICKSTART_GRAPHQL.md** → Este arquivo (guia rápido)

### 🔧 Arquivos Técnicos
- **graphql_schema.go** → Definição dos tipos GraphQL
- **graphql_resolvers.go** → Implementação dos resolvers
- **main.go** → Configuração do handler GraphQL

### 🧪 Testes e Demos
- **test_introspection.ps1** → Script de teste interativo
- **graphql_demo.html** → Demo visual no navegador
- **requests/graphql_queries.http** → Queries HTTP prontas

---

## 🎯 Checklist de Verificação

Marque o que você já testou:

- [ ] ✅ Abri o GraphiQL no navegador
- [ ] ✅ Testei o auto-complete (Ctrl + Space)
- [ ] ✅ Explorei a documentação (botão Docs)
- [ ] ✅ Executei uma query de introspection
- [ ] ✅ Formatei uma query automaticamente (Ctrl + Shift + P)
- [ ] ✅ Executei uma query com múltiplos campos
- [ ] ✅ Testei query com variáveis
- [ ] ✅ Executei o script test_introspection.ps1

---

## 💡 Dicas Pro

### Dica 1: Use Aliases
```graphql
query {
  cliente1: countReceivablesByCustomer(codigo_cliente: "CLI-001") {
    count
  }
  cliente2: countReceivablesByCustomer(codigo_cliente: "CLI-002") {
    count
  }
}
```

### Dica 2: Use Fragmentos
```graphql
fragment ReceivableInfo on Receivable {
  id
  codigo_cliente
  valor_original
}

query {
  getReceivableById(id: "REC-001") {
    ...ReceivableInfo
  }
}
```

### Dica 3: Use Variáveis
```graphql
query GetBalance($cliente: String!) {
  getCustomerBalance(
    codigo_cliente: $cliente
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
  ) {
    saldo_total
  }
}
```

Variáveis (painel Query Variables):
```json
{
  "cliente": "123456"
}
```

---

## 🚀 Próximos Passos

Agora que você conhece GraphQL + Introspection:

1. **Explore** todas as 10 queries disponíveis
2. **Teste** diferentes combinações de campos
3. **Crie** queries customizadas para seu uso
4. **Integre** com frontend (React, Vue, Angular, etc)
5. **Use** ferramentas como Apollo Client ou Relay

---

## 🆘 Problemas Comuns

### GraphiQL não abre?
```powershell
# Verifique se o servidor está rodando
curl http://localhost:8080/health

# Reinicie o servidor
go run .
```

### Auto-complete não funciona?
- Certifique-se de que a query está sintaticamente correta
- Use `Ctrl + Space` logo após `{` ou nome de campo

### Erro ao executar query?
- Verifique o botão "Docs" para ver argumentos obrigatórios
- Use auto-complete para evitar erros de digitação

---

## 📞 Referências Rápidas

| Recurso | URL/Comando |
|---------|-------------|
| GraphiQL | http://localhost:8080/graphql |
| Health Check | http://localhost:8080/health |
| Guia Completo | [GRAPHQL_GUIDE.md](./GRAPHQL_GUIDE.md) |
| Exemplos | [graphql_queries_examples.md](./graphql_queries_examples.md) |
| Script de Teste | `.\test_introspection.ps1` |
| Demo HTML | Abra `graphql_demo.html` no navegador |

---

## ✅ Resumo

Você tem **tudo pronto e funcionando**:

1. ✅ **GraphiQL** - Interface interativa no navegador
2. ✅ **Introspection** - Documentação automática
3. ✅ **10 Queries** - Prontas para uso
4. ✅ **Auto-complete** - Facilitando o desenvolvimento
5. ✅ **Documentação** - Guias completos
6. ✅ **Scripts de Teste** - Para validação

**Basta abrir o navegador em http://localhost:8080/graphql e começar a explorar!** 🚀

---

**Happy Coding! 🎉**
