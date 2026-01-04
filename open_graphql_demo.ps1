# Script de Abertura Rápida - GraphQL Demo
# Abre todos os recursos GraphQL de uma vez

Write-Host ""
Write-Host "🔷 ========================================" -ForegroundColor Cyan
Write-Host "🔷   GraphQL Demo Launcher" -ForegroundColor Cyan
Write-Host "🔷   Data Aggregator - Elasticsearch API" -ForegroundColor Cyan
Write-Host "🔷 ========================================" -ForegroundColor Cyan
Write-Host ""

# Verificar se o servidor está rodando
Write-Host "🔍 Verificando servidor..." -ForegroundColor Yellow

try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 3
    Write-Host "✅ Servidor está rodando!" -ForegroundColor Green
    Write-Host "   Status: $($health.status)" -ForegroundColor Gray
    Write-Host ""
}
catch {
    Write-Host "❌ Servidor NÃO está rodando!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Deseja iniciar o servidor agora? (S/N)" -ForegroundColor Yellow
    $resposta = Read-Host
    
    if ($resposta -eq "S" -or $resposta -eq "s") {
        Write-Host ""
        Write-Host "🚀 Iniciando servidor..." -ForegroundColor Cyan
        Write-Host "   Aguarde o servidor iniciar completamente antes de usar o GraphiQL" -ForegroundColor Yellow
        Write-Host ""
        
        # Iniciar servidor em nova janela
        Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot'; go run ."
        
        Write-Host "⏳ Aguardando servidor iniciar..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
        
        # Tentar conectar novamente
        $tentativas = 0
        $maxTentativas = 10
        
        while ($tentativas -lt $maxTentativas) {
            try {
                $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 2
                Write-Host "✅ Servidor iniciado com sucesso!" -ForegroundColor Green
                Write-Host ""
                break
            }
            catch {
                $tentativas++
                Write-Host "   Tentativa $tentativas/$maxTentativas..." -ForegroundColor Gray
                Start-Sleep -Seconds 2
            }
        }
        
        if ($tentativas -eq $maxTentativas) {
            Write-Host "❌ Não foi possível conectar ao servidor" -ForegroundColor Red
            Write-Host "   Verifique se o servidor está rodando na janela aberta" -ForegroundColor Yellow
            Write-Host ""
            exit 1
        }
    }
    else {
        Write-Host ""
        Write-Host "⚠️  Inicie o servidor manualmente com:" -ForegroundColor Yellow
        Write-Host "   go run ." -ForegroundColor White
        Write-Host ""
        exit 1
    }
}

# Menu de opções
Write-Host "Escolha uma opção:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. 🎨 Abrir GraphiQL no navegador" -ForegroundColor White
Write-Host "2. 📖 Abrir documentação HTML" -ForegroundColor White
Write-Host "3. 📚 Abrir guia completo (Markdown)" -ForegroundColor White
Write-Host "4. 🧪 Executar testes de introspection" -ForegroundColor White
Write-Host "5. 🚀 Abrir TUDO de uma vez!" -ForegroundColor Green
Write-Host "0. ❌ Sair" -ForegroundColor White
Write-Host ""

$opcao = Read-Host "Digite o número da opção"
Write-Host ""

function Open-GraphiQL {
    Write-Host "🎨 Abrindo GraphiQL..." -ForegroundColor Cyan
    Start-Process "http://localhost:8080/graphql"
    Write-Host "✅ GraphiQL aberto no navegador!" -ForegroundColor Green
}

function Open-HTMLDemo {
    Write-Host "📖 Abrindo documentação HTML..." -ForegroundColor Cyan
    $htmlPath = Join-Path $PSScriptRoot "graphql_demo.html"
    Start-Process $htmlPath
    Write-Host "✅ Documentação HTML aberta!" -ForegroundColor Green
}

function Open-MarkdownGuide {
    Write-Host "📚 Abrindo guia completo..." -ForegroundColor Cyan
    $guidePath = Join-Path $PSScriptRoot "GRAPHQL_GUIDE.md"
    Start-Process $guidePath
    Write-Host "✅ Guia aberto no editor padrão!" -ForegroundColor Green
}

function Run-IntrospectionTests {
    Write-Host "🧪 Executando testes de introspection..." -ForegroundColor Cyan
    Write-Host ""
    $testScript = Join-Path $PSScriptRoot "test_introspection.ps1"
    & $testScript
}

switch ($opcao) {
    "1" {
        Open-GraphiQL
    }
    "2" {
        Open-HTMLDemo
    }
    "3" {
        Open-MarkdownGuide
    }
    "4" {
        Run-IntrospectionTests
    }
    "5" {
        Write-Host "🚀 Abrindo TUDO!" -ForegroundColor Green
        Write-Host ""
        
        Open-GraphiQL
        Start-Sleep -Seconds 1
        
        Open-HTMLDemo
        Start-Sleep -Seconds 1
        
        Open-MarkdownGuide
        Start-Sleep -Seconds 1
        
        Write-Host ""
        Write-Host "✅ Tudo aberto!" -ForegroundColor Green
        Write-Host ""
        Write-Host "📋 Recursos abertos:" -ForegroundColor Cyan
        Write-Host "   • GraphiQL (navegador)" -ForegroundColor Gray
        Write-Host "   • Documentação HTML (navegador)" -ForegroundColor Gray
        Write-Host "   • Guia Markdown (editor)" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Deseja executar os testes de introspection? (S/N)" -ForegroundColor Yellow
        $resposta = Read-Host
        
        if ($resposta -eq "S" -or $resposta -eq "s") {
            Write-Host ""
            Run-IntrospectionTests
        }
    }
    "0" {
        Write-Host "👋 Até logo!" -ForegroundColor Yellow
        exit 0
    }
    default {
        Write-Host "❌ Opção inválida!" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
Write-Host ""
Write-Host "💡 Dicas rápidas:" -ForegroundColor Cyan
Write-Host ""
Write-Host "   • No GraphiQL, pressione Ctrl+Space para auto-complete" -ForegroundColor Gray
Write-Host "   • Clique em 'Docs' para ver toda a documentação" -ForegroundColor Gray
Write-Host "   • Pressione Ctrl+Enter para executar queries" -ForegroundColor Gray
Write-Host "   • Use Ctrl+Shift+P para formatar queries" -ForegroundColor Gray
Write-Host ""
Write-Host "📚 Arquivos de referência:" -ForegroundColor Cyan
Write-Host "   • QUICKSTART_GRAPHQL.md - Guia de início rápido" -ForegroundColor Gray
Write-Host "   • GRAPHQL_GUIDE.md - Guia completo" -ForegroundColor Gray
Write-Host "   • graphql_queries_examples.md - Exemplos de queries" -ForegroundColor Gray
Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
Write-Host ""
Write-Host "✅ Pronto para explorar GraphQL!" -ForegroundColor Green
Write-Host ""
