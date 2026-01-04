# Script de Teste - GraphQL Introspection
# Data Aggregator

Write-Host "🔷 GraphQL Introspection Test Script" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

$graphqlUrl = "http://localhost:8080/graphql"

# Função para fazer requisições GraphQL
function Invoke-GraphQLQuery {
    param(
        [string]$Query,
        [string]$TestName
    )
    
    Write-Host "📋 Teste: $TestName" -ForegroundColor Yellow
    Write-Host "Query:" -ForegroundColor Gray
    Write-Host $Query -ForegroundColor DarkGray
    Write-Host ""
    
    $body = @{
        query = $Query
    } | ConvertTo-Json -Depth 10
    
    try {
        $response = Invoke-RestMethod -Uri $graphqlUrl -Method Post -Body $body -ContentType "application/json"
        Write-Host "✅ Resposta:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 10 | Write-Host
    }
    catch {
        Write-Host "❌ Erro: $_" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "-----------------------------------" -ForegroundColor DarkGray
    Write-Host ""
}

# Verificar se o servidor está rodando
Write-Host "🔍 Verificando se o servidor está rodando..." -ForegroundColor White
try {
    $healthCheck = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "✅ Servidor está rodando!" -ForegroundColor Green
    Write-Host "   Status: $($healthCheck.status)" -ForegroundColor Gray
    Write-Host "   Service: $($healthCheck.service)" -ForegroundColor Gray
    Write-Host ""
}
catch {
    Write-Host "❌ Servidor não está rodando!" -ForegroundColor Red
    Write-Host "   Execute 'go run .' no diretório do projeto primeiro." -ForegroundColor Yellow
    Write-Host ""
    exit 1
}

# Menu de testes
Write-Host "Selecione um teste para executar:" -ForegroundColor Cyan
Write-Host "1. Listar todas as queries disponíveis" -ForegroundColor White
Write-Host "2. Ver detalhes do tipo Query" -ForegroundColor White
Write-Host "3. Ver detalhes do tipo Receivable" -ForegroundColor White
Write-Host "4. Ver detalhes do tipo Balance" -ForegroundColor White
Write-Host "5. Ver todos os tipos disponíveis" -ForegroundColor White
Write-Host "6. Ver schema completo" -ForegroundColor White
Write-Host "7. Executar todos os testes" -ForegroundColor White
Write-Host "0. Sair" -ForegroundColor White
Write-Host ""

$choice = Read-Host "Digite o número do teste"

switch ($choice) {
    "1" {
        $query = @"
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
            kind
          }
          defaultValue
        }
      }
    }
  }
}
"@
        Invoke-GraphQLQuery -Query $query -TestName "Listar todas as queries disponíveis"
    }
    
    "2" {
        $query = @"
{
  __type(name: "Query") {
    name
    kind
    fields {
      name
      description
      args {
        name
        description
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
"@
        Invoke-GraphQLQuery -Query $query -TestName "Ver detalhes do tipo Query"
    }
    
    "3" {
        $query = @"
{
  __type(name: "Receivable") {
    name
    kind
    description
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
"@
        Invoke-GraphQLQuery -Query $query -TestName "Ver detalhes do tipo Receivable"
    }
    
    "4" {
        $query = @"
{
  __type(name: "Balance") {
    name
    kind
    description
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
"@
        Invoke-GraphQLQuery -Query $query -TestName "Ver detalhes do tipo Balance"
    }
    
    "5" {
        $query = @"
{
  __schema {
    types {
      name
      kind
      description
    }
  }
}
"@
        Invoke-GraphQLQuery -Query $query -TestName "Ver todos os tipos disponíveis"
    }
    
    "6" {
        $query = @"
{
  __schema {
    queryType {
      name
    }
    types {
      name
      kind
      fields {
        name
        type {
          name
        }
      }
    }
  }
}
"@
        Invoke-GraphQLQuery -Query $query -TestName "Ver schema completo"
    }
    
    "7" {
        Write-Host "🚀 Executando todos os testes..." -ForegroundColor Cyan
        Write-Host ""
        
        # Teste 1: Listar queries
        $query1 = "{ __schema { queryType { fields { name description } } } }"
        Invoke-GraphQLQuery -Query $query1 -TestName "1. Listar queries"
        
        # Teste 2: Tipo Receivable
        $query2 = "{ __type(name: \"Receivable\") { name fields { name type { name } } } }"
        Invoke-GraphQLQuery -Query $query2 -TestName "2. Tipo Receivable"
        
        # Teste 3: Tipo Balance
        $query3 = "{ __type(name: \"Balance\") { name fields { name type { name } } } }"
        Invoke-GraphQLQuery -Query $query3 -TestName "3. Tipo Balance"
        
        # Teste 4: Todos os tipos
        $query4 = "{ __schema { types { name kind } } }"
        Invoke-GraphQLQuery -Query $query4 -TestName "4. Todos os tipos"
        
        # Teste 5: Query real - Index Count
        $query5 = "{ getIndexCount { count } }"
        Invoke-GraphQLQuery -Query $query5 -TestName "5. Query Real - Index Count"
        
        Write-Host "✅ Todos os testes concluídos!" -ForegroundColor Green
    }
    
    "0" {
        Write-Host "👋 Saindo..." -ForegroundColor Yellow
        exit 0
    }
    
    default {
        Write-Host "❌ Opção inválida!" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "✅ Script finalizado!" -ForegroundColor Green
Write-Host ""
Write-Host "💡 Dica: Abra o GraphiQL em http://localhost:8080/graphql" -ForegroundColor Cyan
Write-Host "   para uma experiência interativa completa!" -ForegroundColor Cyan
Write-Host ""
