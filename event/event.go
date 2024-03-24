package event

import (
	"fmt"
	"sync"
)

// const (
// 	AddedToCart          = "AddedToCart"
// 	PurchaseConfirmation = "PurchaseConfirmation"
// )

var EB *eventBus

// type EventBus interface {
// 	Subscribe(topic string) <-chan string
// 	Publish(topic string, payload T)
// 	Close()
// }

type eventBus struct {
	mu     sync.RWMutex
	subs   map[string][]chan string
	closed bool
}

func NewEventBus() *eventBus {
	EB = &eventBus{}
	EB.subs = make(map[string][]chan string)
	EB.closed = false
	return EB
}

// consume
func (eb *eventBus) Subscribe(topic string) <-chan string {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan string, 1) // check buffer size?
	eb.subs[topic] = append(eb.subs[topic], ch)
	return ch
}

// produce
func (eb *eventBus) Publish(topic string, payload string) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return
	}

	for _, ch := range eb.subs[topic] {
		ch <- payload
	}
}

func (eb *eventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if !eb.closed {
		eb.closed = true
		for _, subs := range eb.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
	fmt.Println("EventBus Closed")
}

// type PayloadTypes interface {
// 	AddedToCartPayload | CheckoutPayload
// }

// type AddedToCartPayload struct {
// 	ItemSku  string
// 	CartId   string
// 	ServerTS time.Time
// }

// type CheckoutPayload struct {
// 	CartId   string
// 	ServerTS time.Time
// }
