package event

import (
	"fmt"
	"sync"
)

type Topic string

const (
	TopicAddedToCart       Topic = "AddedToCart"
	TopicCheckout          Topic = "Checkout"
	TopicPurchaseConfirmed Topic = "PurchaseConfirmed"
)

var EB *eventBus

type eventBus struct {
	mu     sync.RWMutex
	subs   map[Topic][]chan string
	closed bool
}

func NewEventBus() *eventBus {
	EB = &eventBus{}
	EB.subs = make(map[Topic][]chan string)
	EB.closed = false
	return EB
}

// consume
func (eb *eventBus) Subscribe(topic Topic) <-chan string {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan string, 1) // check buffer size?
	eb.subs[topic] = append(eb.subs[topic], ch)
	return ch
}

// produce
func (eb *eventBus) Publish(topic Topic, payload string) {
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
}

func NewDefaultHandler(subscriber string, ch <-chan string) {
	for event := range ch {
		fmt.Printf("[%s] got %s\n", subscriber, event)
	}
}
