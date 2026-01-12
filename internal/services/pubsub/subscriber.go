package pubsub

import (
	"project-s/internal/types"
	"sync"
)

type Subscriber struct {
	SubscribeTo  map[*Publisher]map[string]bool //Key: Pub Value: map[ Key: subTopic value: bool ]
	playerAction chan types.PlayerAction
	mu           sync.RWMutex
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		SubscribeTo:  make(map[*Publisher]map[string]bool),
		playerAction: make(chan types.PlayerAction, 5),
	}
}

func (sub *Subscriber) Subscribe(publisher *Publisher, topic string, callback func(actionEvent types.PlayerAction)) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	if publisher.Subscribers[topic] == nil {
		publisher.Subscribers[topic] = make(map[*Subscriber]bool)
	}

	if sub.SubscribeTo[publisher] == nil {
		sub.SubscribeTo[publisher] = make(map[string]bool)
	}

	if _, exist := publisher.Subscribers[topic][sub]; exist {
		return
	}

	publisher.Subscribers[topic][sub] = true
	sub.SubscribeTo[publisher][topic] = true

	go func() {
		for playerAction := range sub.playerAction {
			callback(playerAction)
		}
	}()

}

func (sub *Subscriber) Unsubscribe(topic string, publisher *Publisher) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	delete(publisher.Subscribers[topic], sub)
	delete(sub.SubscribeTo[publisher], topic)

	if len(sub.SubscribeTo[publisher]) == 0 {
		delete(sub.SubscribeTo, publisher)
	}

	if len(sub.SubscribeTo) == 0 {
		close(sub.playerAction)
	}
}
