# go-web-store
A GraphQL API paired with an eventing system

## Getting Started
`go run server.go` 
and navigate to http://localhost:8080/

## Event Design
As implemented, these subscribers are abstract services and their consumption of events results in a print statement
| Publisher         | Topic             | Subscriber              |
| :---:             |    :----:         |         :---:           |
| AddItemToCart()   | AddedToCart       | Analytics               |
| Checkout()        | Checkout          | Analytics               |
| ConfirmPurchase() | PurchaseConfirmed | Analytics               |
| ConfirmPurchase() | PurchaseConfirmed | EmailService            |
| ConfirmPurchase() | PurchaseConfirmed | OrderFullfilmentService |

## Interacting with the API
Users must first use `createItem()` to populate the in memory storage with `Items`. 
Then user can `addItemToCart()`, `checkout()` to generate an `Invoice`, and then `confirmPurchase()` to mark an `Invoice` as paid