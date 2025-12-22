# Script de Validação de Cálculo de Saldo
# Compara 3 métodos: Painless script, Agregações manuais e Bucket Script (performático)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "VALIDAÇÃO DE CÁLCULO DE SALDO" -ForegroundColor Cyan
Write-Host "3 Métodos: Painless | Agregações | Bucket Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$cliente1 = "CLI-10009"
$cliente2 = "CLI-10010"
$cliente3 = "CLI-10011"
$dataInicio = "2025-01-01"
$dataFim = "2026-12-31"

Write-Host "Clientes: $cliente1, $cliente2, $cliente3" -ForegroundColor Yellow
Write-Host "Período: $dataInicio a $dataFim" -ForegroundColor Yellow
Write-Host ""

# Query 1: Buscar soma de valores originais, cancelamentos e negociações (agregações nativas)
Write-Host "1. Executando query com agregações nativas..." -ForegroundColor Green

$querySomas = @"
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "codigo_cliente.keyword": "$cliente1"
          }
        },
        {
          "range": {
            "data_vencimento": {
              "gte": "$dataInicio",
              "lte": "$dataFim"
            }
          }
        }
      ]
    }
  },
  "aggs": {
    "soma_cancelamentos": {
      "nested": {
        "path": "cancelamentos"
      },
      "aggs": {
        "total_cancelado": {
          "sum": {
            "field": "cancelamentos.valor_cancelado"
          }
        }
      }
    },
    "soma_negociacoes": {
      "nested": {
        "path": "negociacoes"
      },
      "aggs": {
        "total_negociado": {
          "sum": {
            "field": "negociacoes.valor_negociado"
          }
        }
      }
    },
    "soma_valores_originais": {
      "sum": {
        "field": "valor_original"
      }
    }
  }
}
"@

$querySomas | Out-File -FilePath temp_somas.json -Encoding UTF8
$startSomas = Get-Date
$resultSomas = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
  -H "Content-Type: application/json" `
  --data-binary "@temp_somas.json" 2>$null | ConvertFrom-Json
$endSomas = Get-Date
$tempoSomas = ($endSomas - $startSomas).TotalMilliseconds

$valorOriginalTotal = $resultSomas.aggregations.soma_valores_originais.value
$cancelamentoTotal = $resultSomas.aggregations.soma_cancelamentos.total_cancelado.value
$negociacaoTotal = $resultSomas.aggregations.soma_negociacoes.total_negociado.value
$saldoCalculado = $valorOriginalTotal - $cancelamentoTotal - $negociacaoTotal

Write-Host "   Valores Originais: R$ $([math]::Round($valorOriginalTotal, 2))" -ForegroundColor White
Write-Host "   Cancelamentos:     R$ $([math]::Round($cancelamentoTotal, 2))" -ForegroundColor White
Write-Host "   Negociações:       R$ $([math]::Round($negociacaoTotal, 2))" -ForegroundColor White
Write-Host "   Saldo Calculado:   R$ $([math]::Round($saldoCalculado, 2))" -ForegroundColor Cyan
Write-Host "   Tempo:             $([math]::Round($tempoSomas, 2))ms" -ForegroundColor Gray
Write-Host ""

# Query 2: Buscar saldo via Painless script
Write-Host "2. Executando query com Painless script..." -ForegroundColor Green

$queryPainless = @"
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "codigo_cliente.keyword": "$cliente2"
          }
        },
        {
          "range": {
            "data_vencimento": {
              "gte": "$dataInicio",
              "lte": "$dataFim"
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

$queryPainless | Out-File -FilePath temp_painless.json -Encoding UTF8
$startPainless = Get-Date
$resultPainless = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
  -H "Content-Type: application/json" `
  --data-binary "@temp_painless.json" 2>$null | ConvertFrom-Json
$endPainless = Get-Date
$tempoPainless = ($endPainless - $startPainless).TotalMilliseconds

$saldoPainless = $resultPainless.aggregations.saldo_total.value

Write-Host "   Saldo (Painless):  R$ $([math]::Round($saldoPainless, 2))" -ForegroundColor Cyan
Write-Host "   Tempo:             $([math]::Round($tempoPainless, 2))ms" -ForegroundColor Gray
Write-Host ""

# Query 3: Buscar saldo via Bucket Script (performático)
Write-Host "3. Executando query com Bucket Script (performático)..." -ForegroundColor Green

$queryBucketScript = @"
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "codigo_cliente.keyword": "$cliente3"
          }
        },
        {
          "range": {
            "data_vencimento": {
              "gte": "$dataInicio",
              "lte": "$dataFim"
            }
          }
        }
      ]
    }
  },
  "aggs": {
    "soma_valores_originais": {
      "sum": {
        "field": "valor_original"
      }
    },
    "soma_cancelamentos": {
      "nested": {
        "path": "cancelamentos"
      },
      "aggs": {
        "total_cancelado": {
          "sum": {
            "field": "cancelamentos.valor_cancelado"
          }
        }
      }
    },
    "soma_negociacoes": {
      "nested": {
        "path": "negociacoes"
      },
      "aggs": {
        "total_negociado": {
          "sum": {
            "field": "negociacoes.valor_negociado"
          }
        }
      }
    },
    "saldo_disponivel": {
      "bucket_script": {
        "buckets_path": {
          "valores": "soma_valores_originais",
          "cancelamentos": "soma_cancelamentos>total_cancelado",
          "negociacoes": "soma_negociacoes>total_negociado"
        },
        "script": "Math.round((params.valores - params.cancelamentos - params.negociacoes) * 100) / 100"
      }
    }
  }
}
"@

$queryBucketScript | Out-File -FilePath temp_bucket.json -Encoding UTF8
$startBucket = Get-Date
$resultBucket = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search" `
  -H "Content-Type: application/json" `
  --data-binary "@temp_bucket.json" 2>$null | ConvertFrom-Json
$endBucket = Get-Date
$tempoBucket = ($endBucket - $startBucket).TotalMilliseconds

$saldoBucket = $resultBucket.aggregations.saldo_disponivel.value

Write-Host "   Saldo (Bucket Script): R$ $([math]::Round($saldoBucket, 2))" -ForegroundColor Cyan
Write-Host "   Tempo:                 $([math]::Round($tempoBucket, 2))ms" -ForegroundColor Gray
Write-Host ""

# Validação
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "RESULTADO DA VALIDAÇÃO" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$diferencaPainless = [math]::Abs($saldoCalculado - $saldoPainless)
$diferencaBucket = [math]::Abs($saldoCalculado - $saldoBucket)
$diferencaPainlessBucket = [math]::Abs($saldoPainless - $saldoBucket)

Write-Host "Método 1 - Agregações Manuais:  R$ $([math]::Round($saldoCalculado, 2))" -ForegroundColor Yellow
Write-Host "Método 2 - Painless Script:     R$ $([math]::Round($saldoPainless, 2))" -ForegroundColor Yellow
Write-Host "Método 3 - Bucket Script:       R$ $([math]::Round($saldoBucket, 2))" -ForegroundColor Yellow
Write-Host ""

$todosIguais = ($diferencaPainless -lt 0.01) -and ($diferencaBucket -lt 0.01) -and ($diferencaPainlessBucket -lt 0.01)

if ($todosIguais) {
    Write-Host "✅ VALIDAÇÃO PASSOU: Todos os métodos retornam valores equivalentes!" -ForegroundColor Green
} else {
    Write-Host "❌ VALIDAÇÃO FALHOU: Diferenças detectadas!" -ForegroundColor Red
    Write-Host "   Agregações vs Painless: R$ $([math]::Round($diferencaPainless, 2))" -ForegroundColor $(if ($diferencaPainless -lt 0.01) { "Green" } else { "Red" })
    Write-Host "   Agregações vs Bucket:   R$ $([math]::Round($diferencaBucket, 2))" -ForegroundColor $(if ($diferencaBucket -lt 0.01) { "Green" } else { "Red" })
    Write-Host "   Painless vs Bucket:     R$ $([math]::Round($diferencaPainlessBucket, 2))" -ForegroundColor $(if ($diferencaPainlessBucket -lt 0.01) { "Green" } else { "Red" })
}

Write-Host ""
Write-Host "Performance (menor é melhor):" -ForegroundColor Yellow
Write-Host "  1. Agregações Manuais:  $([math]::Round($tempoSomas, 2))ms" -ForegroundColor $(if ($tempoSomas -le $tempoPainless -and $tempoSomas -le $tempoBucket) { "Green" } else { "White" })
Write-Host "  2. Painless Script:     $([math]::Round($tempoPainless, 2))ms" -ForegroundColor $(if ($tempoPainless -le $tempoSomas -and $tempoPainless -le $tempoBucket) { "Green" } else { "White" })
Write-Host "  3. Bucket Script:       $([math]::Round($tempoBucket, 2))ms" -ForegroundColor $(if ($tempoBucket -le $tempoSomas -and $tempoBucket -le $tempoPainless) { "Green" } else { "White" })

Write-Host ""

# Determinar o método mais rápido
$metodos = @(
    @{Nome="Agregações Manuais"; Tempo=$tempoSomas},
    @{Nome="Painless Script"; Tempo=$tempoPainless},
    @{Nome="Bucket Script"; Tempo=$tempoBucket}
)
$maisRapido = $metodos | Sort-Object Tempo | Select-Object -First 1
$maisLento = $metodos | Sort-Object Tempo -Descending | Select-Object -First 1

$ganho = (($maisLento.Tempo - $maisRapido.Tempo) / $maisLento.Tempo) * 100

Write-Host "🏆 Método mais rápido: $($maisRapido.Nome)" -ForegroundColor Cyan
Write-Host "   Ganho sobre o mais lento: $([math]::Round($ganho, 2))%" -ForegroundColor Green

Write-Host ""
Write-Host "💡 Recomendação:" -ForegroundColor Yellow
if ($tempoBucket -le $tempoPainless * 0.8) {
    Write-Host "   Use Bucket Script para melhor performance e simplicidade!" -ForegroundColor Green
} elseif ($tempoSomas -le $tempoPainless * 0.8) {
    Write-Host "   Agregações Manuais são eficientes, mas Bucket Script retorna saldo direto!" -ForegroundColor Cyan
} else {
    Write-Host "   Todos os métodos têm performance similar." -ForegroundColor White
}

Write-Host ""

# Limpar arquivos temporários
Remove-Item temp_somas.json -ErrorAction SilentlyContinue
Remove-Item temp_painless.json -ErrorAction SilentlyContinue
Remove-Item temp_bucket.json -ErrorAction SilentlyContinue

Write-Host "========================================" -ForegroundColor Cyan
