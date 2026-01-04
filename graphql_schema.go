package main

import (
	"github.com/graphql-go/graphql"
)

// GraphQL Types
var receivableType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Receivable",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"id_recebivel": &graphql.Field{
			Type: graphql.String,
		},
		"codigo_cliente": &graphql.Field{
			Type: graphql.String,
		},
		"valor_original": &graphql.Field{
			Type: graphql.Float,
		},
		"data_vencimento": &graphql.Field{
			Type: graphql.String,
		},
		"id_pagamento": &graphql.Field{
			Type: graphql.String,
		},
		"saldo_disponivel": &graphql.Field{
			Type: graphql.Float,
		},
		"cancelamentos": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Cancelamento",
				Fields: graphql.Fields{
					"data_cancelamento": &graphql.Field{Type: graphql.String},
					"valor_cancelado":   &graphql.Field{Type: graphql.Float},
					"motivo":            &graphql.Field{Type: graphql.String},
				},
			})),
		},
		"negociacoes": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Negociacao",
				Fields: graphql.Fields{
					"data_negociacao": &graphql.Field{Type: graphql.String},
					"valor_negociado": &graphql.Field{Type: graphql.Float},
					"tipo_negociacao": &graphql.Field{Type: graphql.String},
					"observacao":      &graphql.Field{Type: graphql.String},
				},
			})),
		},
	},
})

var balanceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Balance",
	Fields: graphql.Fields{
		"codigo_cliente": &graphql.Field{
			Type: graphql.String,
		},
		"periodo": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "Period",
				Fields: graphql.Fields{
					"inicio": &graphql.Field{Type: graphql.String},
					"fim":    &graphql.Field{Type: graphql.String},
				},
			}),
		},
		"total_recebiveis": &graphql.Field{
			Type: graphql.Int,
		},
		"saldo_total": &graphql.Field{
			Type: graphql.Float,
		},
		"saldo_formatado": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var customerStatsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CustomerStats",
	Fields: graphql.Fields{
		"codigo_cliente": &graphql.Field{
			Type: graphql.String,
		},
		"total_recebiveis": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var searchResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SearchResult",
	Fields: graphql.Fields{
		"total": &graphql.Field{
			Type: graphql.Int,
		},
		"receivables": &graphql.Field{
			Type: graphql.NewList(receivableType),
		},
	},
})

var countResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CountResult",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
