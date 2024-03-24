package graph

import "github.com/cpaules/go-web-store/graph/model"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	items    map[string]*model.Item
	cart     *model.Cart
	invoices []*model.Invoice
}

func NewDefaultResolver() *Resolver {
	return &Resolver{
		items: make(map[string]*model.Item),
		cart: &model.Cart{
			ID:    "1",
			Items: []*model.Item{},
		},
		invoices: []*model.Invoice{},
	}
}
