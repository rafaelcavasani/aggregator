# Script de Teste de Todas as Queries GraphQL
# Com 10 milhões de registros

Write-Host ""
Write-Host "🔷 ================================================" -ForegroundColor Cyan
Write-Host "🔷   TESTE DE QUERIES GRAPHQL" -ForegroundColor Cyan
Write-Host "🔷   10 Milhões de Registros no Elasticsearch" -ForegroundColor Cyan
Write-Host "🔷 ================================================" -ForegroundColor Cyan
Write-Host ""

$graphqlUrl = "http://localhost:8080/graphql"

# Verificar se o servidor está rodando
Write-Host "🔍 Verificando servidor..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "✅ Servidor está rodando!" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Servidor não está rodando!" -ForegroundColor Red
    Write-Host "   Execute 'go run .' primeiro" -ForegroundColor Yellow
    exit 1
}

# Definir queries para teste
$queries = @(
    @{
        name = "1. getIndexCount"
        description = "Contar total de documentos"
        query = @'
{
  getIndexCount {
    count
  }
}
'@
    },
    @{
        name = "2. getTopCustomer"
        description = "Cliente com mais registros"
        query = @'
{
  getTopCustomer {
    codigo_cliente
    total_recebiveis
  }
}
'@
    },
    @{
        name = "3. getAllReceivables"
        description = "Buscar primeiros 5 recebíveis"
        query = @'
{
  getAllReceivables(size: 5) {
    total
    receivables {
      id_recebivel
      codigo_cliente
      valor_original
    }
  }
}
'@
    },
    @{
        name = "4. getReceivableById"
        description = "Buscar por ID específico"
        query = @'
{
  getReceivableById(id: "47e80a94-e661-4de9-9f46-9faabdac709b") {
    id_recebivel
    codigo_cliente
    valor_original
    data_vencimento
  }
}
'@
    },
    @{
        name = "5. countReceivablesByCustomer"
        description = "Contar recebíveis de um cliente"
        query = @'
{
  countReceivablesByCustomer(codigo_cliente: "CLI-10001") {
    count
  }
}
'@
    },
    @{
        name = "6. getReceivablesByCustomerAndDueDate"
        description = "Buscar por cliente e data"
        query = @'
{
  getReceivablesByCustomerAndDueDate(
    codigo_cliente: "CLI-10001"
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
    size: 3
  ) {
    total
    receivables {
      id_recebivel
      valor_original
      data_vencimento
    }
  }
}
'@
    },
    @{
        name = "7. countReceivablesGroupByCustomer"
        description = "Contagem agrupada por cliente"
        query = @'
{
  countReceivablesGroupByCustomer(
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
  ) {
    codigo_cliente
    total_recebiveis
  }
}
'@
    },
    @{
        name = "8. getCustomerBalance"
        description = "Saldo do cliente"
        query = @'
{
  getCustomerBalance(
    codigo_cliente: "CLI-10001"
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
  ) {
    codigo_cliente
    saldo_total
    saldo_formatado
    total_recebiveis
  }
}
'@
    },
    @{
        name = "9. getReceivablesByBalanceAvailable"
        description = "Por saldo disponível"
        query = @'
{
  getReceivablesByBalanceAvailable(
    codigo_cliente: "CLI-10001"
    min_balance: 500.0
    size: 3
  ) {
    total
    receivables {
      id_recebivel
      valor_original
    }
  }
}
'@
    },
    @{
        name = "10. getReceivableBalanceById"
        description = "Saldo calculado por ID"
        query = @'
{
  getReceivableBalanceById(id: "47e80a94-e661-4de9-9f46-9faabdac709b") {
    id_recebivel
    codigo_cliente
    valor_original
    saldo_disponivel
  }
}
'@
    }
)

$successCount = 0
$errorCount = 0
$results = @()

Write-Host "🚀 Iniciando testes..`n" -ForegroundColor Cyan

foreach ($test in $queries) {
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
    Write-Host "📋 $($test.name): $($test.description)" -ForegroundColor Yellow
    
    $body = @{
        query = $test.query
    } | ConvertTo-Json -Depth 10
    
    $startTime = Get-Date
    
    try {
        $response = Invoke-RestMethod -Uri $graphqlUrl -Method Post -Body $body -ContentType "application/json" -TimeoutSec 60
        $endTime = Get-Date
        $duration = ($endTime - $startTime).TotalMilliseconds
        
        if ($response.errors) {
            Write-Host "❌ ERRO!" -ForegroundColor Red
            Write-Host "   Tempo: $([math]::Round($duration, 2))ms" -ForegroundColor Gray
            Write-Host "   Erro: $($response.errors[0].message)" -ForegroundColor Red
            $errorCount++
            
            $results += @{
                name = $test.name
                status = "ERRO"
                time = $duration
                error = $response.errors[0].message
            }
        } else {
            Write-Host "✅ SUCESSO!" -ForegroundColor Green
            Write-Host "   Tempo: $([math]::Round($duration, 2))ms" -ForegroundColor Gray
            $successCount++
            
            $results += @{
                name = $test.name
                status = "SUCESSO"
                time = $duration
            }
            
            # Exibir resultado resumido
            $resultData = $response.data | ConvertTo-Json -Depth 3 -Compress
            if ($resultData.Length -gt 200) {
                Write-Host "   Resultado: $($resultData.Substring(0, 200))..." -ForegroundColor DarkGray
            } else {
                Write-Host "   Resultado: $resultData" -ForegroundColor DarkGray
            }
        }
    }
    catch {
        Write-Host "❌ FALHA NA REQUISIÇÃO!" -ForegroundColor Red
        Write-Host "   Erro: $_" -ForegroundColor Red
        $errorCount++
        
        $results += @{
            name = $test.name
            status = "FALHA"
            error = $_.Exception.Message
        }
    }
    
    Write-Host ""
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
Write-Host ""
Write-Host "📊 RESUMO DOS TESTES" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
Write-Host ""
Write-Host "   Total de queries testadas: $($queries.Count)" -ForegroundColor White
Write-Host "   ✅ Sucesso: $successCount" -ForegroundColor Green
Write-Host "   ❌ Erros: $errorCount" -ForegroundColor Red
Write-Host ""

if ($successCount -gt 0) {
    Write-Host "⏱️  TEMPOS DE EXECUÇÃO:" -ForegroundColor Cyan
    $successResults = $results | Where-Object { $_.status -eq "SUCESSO" }
    foreach ($r in $successResults) {
        $timeColor = if ($r.time -lt 100) { "Green" } elseif ($r.time -lt 1000) { "Yellow" } else { "Red" }
        Write-Host "   $($r.name): $([math]::Round($r.time, 2))ms" -ForegroundColor $timeColor
    }
    Write-Host ""
}

if ($errorCount -gt 0) {
    Write-Host "⚠️  ERROS ENCONTRADOS:" -ForegroundColor Red
    $errorResults = $results | Where-Object { $_.status -ne "SUCESSO" }
    foreach ($r in $errorResults) {
        Write-Host "   $($r.name): $($r.error)" -ForegroundColor Red
    }
    Write-Host ""
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
Write-Host ""

if ($errorCount -eq 0) {
    Write-Host "🎉 TODOS OS TESTES PASSARAM COM SUCESSO!" -ForegroundColor Green
    Write-Host "   Todas as $($queries.Count) queries GraphQL estão funcionando perfeitamente!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Alguns testes falharam. Verifique os erros acima." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "💡 Dica: Abra http://localhost:8080/graphql para testar interativamente" -ForegroundColor Cyan
Write-Host ""
