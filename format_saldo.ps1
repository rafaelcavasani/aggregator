$result = curl -X GET "http://localhost:9200/ciclo_vida_recebivel/_search?pretty" `
  -H "Content-Type: application/json" `
  --data-binary "@query_saldo_cliente.json" `
  2>$null | ConvertFrom-Json

$saldo = [decimal]$result.aggregations.saldo_total.value
$count = $result.hits.total.value

Write-Host "================================" -ForegroundColor Cyan
Write-Host "Cliente: CLI-10001" -ForegroundColor Yellow
Write-Host "Período: 2025-01-01 a 2025-02-28" -ForegroundColor Yellow
Write-Host "Recebíveis: $count" -ForegroundColor Green
Write-Host "Saldo Total: R$ $($saldo.ToString('N2'))" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Cyan
