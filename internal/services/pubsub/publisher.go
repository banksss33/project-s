package pubsub

import (
	"project-s/internal/types"
	"sync"
)

type Publisher struct {
	Subscribers map[string]map[*Subscriber]bool
	mu          sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		Subscribers: make(map[string]map[*Subscriber]bool),
	}
}

func (p *Publisher) Notify(topic string, actionEvent types.ActionEvent) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for subscriber := range p.Subscribers[topic] {
		subscriber.Event <- actionEvent
	}
}
