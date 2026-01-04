# GraphQL Queries Examples
# Use these queries in GraphiQL at http://localhost:8080/graphql
# or send as POST requests to the /graphql endpoint

## 1. Get All Receivables (limit 10)
# Corresponds to: get_all_limit_10.http
query GetAllReceivables {
  getAllReceivables(size: 10) {
    total
    receivables {
      id
      id_recebivel
      codigo_cliente
      valor_original
      data_vencimento
    }
  }
}

## 2. Get Customer Balance by Date Range
# Corresponds to: get_customer_balance_by_date.http
query GetCustomerBalance {
  getCustomerBalance(
    codigo_cliente: "123456"
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
  ) {
    codigo_cliente
    periodo {
      inicio
      fim
    }
    total_recebiveis
    saldo_total
    saldo_formatado
  }
}

## 3. Get Receivable by ID
# Corresponds to: get_document_by_id.http
query GetReceivableById {
  getReceivableById(id: "REC-00001") {
    id
    id_recebivel
    codigo_cliente
    valor_original
    data_vencimento
    id_pagamento
    cancelamentos {
      data_cancelamento
      valor_cancelado
      motivo
    }
    negociacoes {
      data_negociacao
      valor_negociado
      tipo_negociacao
      observacao
    }
  }
}

## 4. Get Receivable Balance by ID with calculated balance
# Corresponds to: get_receivable_balance_by_id.http
query GetReceivableBalanceById {
  getReceivableBalanceById(id: "a7f3c8e2-9b4d-4f6a-8c2e-5d7b9f1a3e6c") {
    id
    id_recebivel
    codigo_cliente
    valor_original
    saldo_disponivel
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

## 5. Get Receivables by Balance Available
# Corresponds to: get_receivables_by_balance_available.http
query GetReceivablesByBalanceAvailable {
  getReceivablesByBalanceAvailable(
    codigo_cliente: "CLI-10001"
    min_balance: 500.00
    from: 0
    size: 50
  ) {
    total
    receivables {
      id
      id_recebivel
      codigo_cliente
      valor_original
      data_vencimento
    }
  }
}

## 6. Get Receivables by Customer and Due Date
# Corresponds to: get_receivables_by_customer_and_due_date.http
query GetReceivablesByCustomerAndDueDate {
  getReceivablesByCustomerAndDueDate(
    codigo_cliente: "123456"
    data_inicio: "2025-12-01"
    data_fim: "2026-12-31"
    from: 0
    size: 10
  ) {
    total
    receivables {
      id
      id_recebivel
      codigo_cliente
      valor_original
      data_vencimento
    }
  }
}

## 7. Count Receivables by Customer
# Corresponds to: get_receivables_count_by_customer.http
query CountReceivablesByCustomer {
  countReceivablesByCustomer(codigo_cliente: "CLI-10008") {
    count
  }
}

## 8. Count Receivables Grouped by Customer
# Corresponds to: get_receivables_count_group_by_customer.http
query CountReceivablesGroupByCustomer {
  countReceivablesGroupByCustomer(
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
  ) {
    codigo_cliente
    total_recebiveis
  }
}

## 9. Get Top Customer
# Corresponds to: get_top_cliente.http
query GetTopCustomer {
  getTopCustomer {
    codigo_cliente
    total_recebiveis
  }
}

## 10. Get Index Count
# Corresponds to: get_index_count.http
query GetIndexCount {
  getIndexCount {
    count
  }
}

## Combined Query Example - Multiple queries in one request
query CombinedQuery {
  indexCount: getIndexCount {
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

## Pagination Example
query PaginatedReceivables($page: Int!, $size: Int!) {
  getReceivablesByCustomerAndDueDate(
    codigo_cliente: "CLI-10001"
    data_inicio: "2025-01-01"
    data_fim: "2026-12-31"
    from: $page
    size: $size
  ) {
    total
    receivables {
      id_recebivel
      codigo_cliente
      valor_original
      data_vencimento
    }
  }
}

# Variables for pagination:
# {
#   "page": 0,
#   "size": 20
# }
