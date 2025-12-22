# Script de Teste de Performance - Top 10 Clientes
# Executa 3 queries em cada um dos 10 clientes com mais recebíveis

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "TESTE DE PERFORMANCE - TOP 10 CLIENTES" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. Buscar top 10 clientes
Write-Host "1. Identificando top 10 clientes..." -ForegroundColor Yellow

$queryTopClientes = @'
{
  "size": 0,
  "aggs": {
    "clientes": {
      "terms": {
        "field": "codigo_cliente.keyword",
        "size": 10,
        "order": {
          "_count": "desc"
        }
      }
    }
  }
}
'@

$queryTopClientes | Out-File -FilePath temp_top10.json -Encoding UTF8
$resultTop10 = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
  -H "Content-Type: application/json" `
  --data-binary "@temp_top10.json" 2>$null | ConvertFrom-Json

$top10Clientes = $resultTop10.aggregations.clientes.buckets

Write-Host "Top 10 clientes identificados:" -ForegroundColor Green
$top10Clientes | ForEach-Object { Write-Host "  - $($_.key): $($_.doc_count) documentos" }
Write-Host ""

# Resultados
$resultados = @()

# 2. Executar testes para cada cliente
foreach ($cliente in $top10Clientes) {
    $codigoCliente = $cliente.key
    $totalDocs = $cliente.doc_count
    
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Cliente: $codigoCliente ($totalDocs docs)" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    
    # Teste 1: Count
    Write-Host "  Teste 1: Contagem..." -ForegroundColor Yellow
    $queryCount = @"
{
  "query": {
    "term": {
      "codigo_cliente.keyword": "$codigoCliente"
    }
  }
}
"@
    
    $queryCount | Out-File -FilePath temp_count.json -Encoding UTF8
    $start = Get-Date
    $resultCount = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_count" `
      -H "Content-Type: application/json" `
      --data-binary "@temp_count.json" 2>$null | ConvertFrom-Json
    $end = Get-Date
    $tempoCount = ($end - $start).TotalMilliseconds
    
    Write-Host "    Count: $($resultCount.count) | Tempo: $([math]::Round($tempoCount, 2))ms" -ForegroundColor Green
    
    # Teste 2: Saldo Total (Painless)
    Write-Host "  Teste 2: Cálculo de saldo (Painless)..." -ForegroundColor Yellow
    $querySaldo = @"
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "codigo_cliente.keyword": "$codigoCliente"
          }
        },
        {
          "range": {
            "data_vencimento": {
              "gte": "2025-01-01",
              "lte": "2026-12-31"
            }
          }
        }
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
"@
    
    $querySaldo | Out-File -FilePath temp_saldo.json -Encoding UTF8
    $start = Get-Date
    $resultSaldo = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
      -H "Content-Type: application/json" `
      --data-binary "@temp_saldo.json" 2>$null | ConvertFrom-Json
    $end = Get-Date
    $tempoSaldo = ($end - $start).TotalMilliseconds
    $saldo = $resultSaldo.aggregations.saldo_total.value
    
    Write-Host "    Saldo: R$ $([math]::Round($saldo, 2)) | Tempo: $([math]::Round($tempoSaldo, 2))ms" -ForegroundColor Green
    
    # Teste 3: Search com ordenação
    Write-Host "  Teste 3: Busca ordenada (100 docs)..." -ForegroundColor Yellow
    $querySearch = @"
{
  "query": {
    "bool": {
      "must": [
        {"term": {"codigo_cliente.keyword": "$codigoCliente"}},
        {"range": {"data_vencimento": {"gte": "2025-01-01", "lte": "2026-12-31"}}}
      ]
    }
  },
  "sort": [{"data_vencimento": {"order": "asc"}}],
  "size": 100
}
"@
    
    $querySearch | Out-File -FilePath temp_search.json -Encoding UTF8
    $start = Get-Date
    $resultSearch = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
      -H "Content-Type: application/json" `
      --data-binary "@temp_search.json" 2>$null | ConvertFrom-Json
    $end = Get-Date
    $tempoSearch = ($end - $start).TotalMilliseconds
    $hitsTotal = $resultSearch.hits.total.value
    $hitsRetornados = $resultSearch.hits.hits.Count
    
    Write-Host "    Hits: $hitsTotal | Retornados: $hitsRetornados | Tempo: $([math]::Round($tempoSearch, 2))ms" -ForegroundColor Green
    Write-Host ""
    
    # Armazenar resultados
    $resultados += [PSCustomObject]@{
        Cliente = $codigoCliente
        TotalDocs = $totalDocs
        CountMs = [math]::Round($tempoCount, 2)
        SaldoMs = [math]::Round($tempoSaldo, 2)
        SearchMs = [math]::Round($tempoSearch, 2)
        Saldo = [math]::Round($saldo, 2)
    }
}

# Limpar arquivos temporários
Remove-Item temp_top10.json -ErrorAction SilentlyContinue
Remove-Item temp_count.json -ErrorAction SilentlyContinue
Remove-Item temp_saldo.json -ErrorAction SilentlyContinue
Remove-Item temp_search.json -ErrorAction SilentlyContinue

# Exibir resumo
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "RESUMO DOS RESULTADOS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$resultados | Format-Table -AutoSize

Write-Host ""
Write-Host "Estatísticas:" -ForegroundColor Yellow
Write-Host "  Count médio: $([math]::Round(($resultados.CountMs | Measure-Object -Average).Average, 2))ms" -ForegroundColor Green
Write-Host "  Saldo médio: $([math]::Round(($resultados.SaldoMs | Measure-Object -Average).Average, 2))ms" -ForegroundColor Green
Write-Host "  Search médio: $([math]::Round(($resultados.SearchMs | Measure-Object -Average).Average, 2))ms" -ForegroundColor Green
Write-Host ""

# Exportar para CSV
$csvPath = "performance_test_results.csv"
$resultados | Export-Csv -Path $csvPath -NoTypeInformation -Encoding UTF8
Write-Host "Resultados exportados para: $csvPath" -ForegroundColor Cyan
