package main

import (
	"fmt"
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
	// add subscribers
	analytics1 := eb.Subscribe("AddedToCart")
	analytics2 := eb.Subscribe("Checkout")
	analytics3 := eb.Subscribe("PurchaseConfirmed")
	emailService := eb.Subscribe("PurchaseConfirmed")
	orderFullfilmentService := eb.Subscribe("PurchaseConfirmed")

	listener := func(subscriber string, ch <-chan string) {
		for event := range ch {
			fmt.Printf("[%s] got %s\n", subscriber, event)
		}
	}

	go listener("Analytics", analytics1)
	go listener("Analytics", analytics2)
	go listener("Analytics", analytics3)
	go listener("EmailService", emailService)
	go listener("OrderFullfilmentService", orderFullfilmentService)

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
