package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cpaules/go-web-store/event"
	"github.com/cpaules/go-web-store/graph/model"
)

// CreateItem is the resolver for the createItem field.
func (r *mutationResolver) CreateItem(ctx context.Context, item model.NewItem) (*model.Item, error) {
	newItem := &model.Item{
		Id:    len(r.items) + 1,
		Price: item.Price,
		Sku:   item.Sku,
	}
	r.items[item.Sku] = newItem
	return newItem, nil
}

// AddItemToCart is the resolver for the addItemToCart field.
func (r *mutationResolver) AddItemToCart(ctx context.Context, sku string) (*model.Cart, error) {
	r.cart.Items = append(r.cart.Items, r.items[sku])
	message := fmt.Sprintf("AddedToCart -> ItemSku: %s, CartId: %s, ServerTS: %s", sku, r.cart.ID, time.Now().Format(time.RFC3339))
	// fmt.Printf("Publishing @%s: %s\n", "AddedToCart", message)
	event.EB.Publish("AddedToCart", message)
	return r.cart, nil
}

// Checkout is the resolver for the checkout field.
func (r *mutationResolver) Checkout(ctx context.Context, cartID string) (*model.Invoice, error) {
	if len(r.cart.Items) == 0 {
		return &model.Invoice{}, errors.New("Cannot checkout with an empty cart!")
	}
	state := model.InvoiceStatePending
	newInvoice := &model.Invoice{
		ID:     strconv.Itoa(len(r.invoices) + 1),
		Items:  r.cart.Items,
		Total:  r.cart.Total(),
		Status: &state,
	}
	r.invoices = append(r.invoices, newInvoice)
	message := fmt.Sprintf("Checkout -> CartId: %s, ServerTS: %s", r.cart.ID, time.Now().Format(time.RFC3339))
	event.EB.Publish("Checkout", message)
	return newInvoice, nil
}

// ConfirmPurchase is the resolver for the confirmPurchase field.
func (r *mutationResolver) ConfirmPurchase(ctx context.Context, invoiceID string, confirm bool) (string, error) {
	invoice, err := r.Query().Invoice(ctx, invoiceID)
	if err != nil {
		return "", err
	}
	r.cart.Items = nil
	// user is declining the invoice, set status to canceled
	if !confirm {
		cancled := model.InvoiceStateCanceled
		invoice.Status = &cancled
		return "Purchase successfully canceled", nil
	}

	paid := model.InvoiceStatePaid
	invoice.Status = &paid

	message := fmt.Sprintf("PurchaseConfirmed -> invoiceID: %s, ServerTS: %s", invoiceID, time.Now().Format(time.RFC3339))
	event.EB.Publish("PurchaseConfirmed", message)
	return "Purchase successfully confirmed", nil
}

// Items is the resolver for the items field.
func (r *queryResolver) Items(ctx context.Context) ([]*model.Item, error) {
	var items []*model.Item
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}

// Cart is the resolver for the cart field.
func (r *queryResolver) Cart(ctx context.Context) (*model.Cart, error) {
	return r.cart, nil
}

// Invoices is the resolver for the invoices field.
func (r *queryResolver) Invoices(ctx context.Context) ([]*model.Invoice, error) {
	return r.invoices, nil
}

// Invoice is the resolver for the invoice field.
func (r *queryResolver) Invoice(ctx context.Context, invoiceID string) (*model.Invoice, error) {
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		return nil, err
	}
	if id > len(r.invoices) {
		return nil, errors.New("invalid invoiceID")
	}
	return r.invoices[id-1], nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
