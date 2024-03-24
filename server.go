package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cpaules/go-web-store/event"
	"github.com/cpaules/go-web-store/graph"
)

const defaultPort = "8080"

func main() {
	eb := event.NewEventBus()

	analytics := map[string]<-chan string{
		"AddedToCart":       eb.Subscribe(event.TopicAddedToCart),
		"Checkout":          eb.Subscribe(event.TopicCheckout),
		"PurchaseConfirmed": eb.Subscribe(event.TopicPurchaseConfirmed),
	}
	emailService := eb.Subscribe(event.TopicPurchaseConfirmed)
	orderFullfilmentService := eb.Subscribe(event.TopicPurchaseConfirmed)

	go event.NewDefaultHandler("Analytics", analytics["AddedToCart"])
	go event.NewDefaultHandler("Analytics", analytics["Checkout"])
	go event.NewDefaultHandler("Analytics", analytics["PurchaseConfirmed"])
	go event.NewDefaultHandler("EmailService", emailService)
	go event.NewDefaultHandler("OrderFullfilmentService", orderFullfilmentService)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewDefaultResolver()}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	// Handle OS interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until an OS interrupt signal is received.
	<-c
	eb.Close()
}
