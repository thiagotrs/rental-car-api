package broker

import "sync"

type pubSub struct {
	mu     sync.RWMutex
	subs   map[string][]chan interface{}
	closed bool
}

func NewPubSub() *pubSub {
	return &pubSub{subs: make(map[string][]chan interface{})}
}

func (ps *pubSub) Subscribe(topic string) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *pubSub) Publish(topic string, data interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		go func(ch chan interface{}) {
			ch <- data
		}(ch)
	}
}

func (ps *pubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
