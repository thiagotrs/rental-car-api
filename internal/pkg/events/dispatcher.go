package events

import "errors"

var ErrNoneHandler = errors.New("no handler for this event")

type eventDispatcher struct {
	EventsHandlers map[string][]EventHandler
}

func NewEventDispatcher() *eventDispatcher {
	return &eventDispatcher{make(map[string][]EventHandler)}
}

func (d eventDispatcher) Register(h EventHandler, eventName string) {
	d.EventsHandlers[eventName] = append(d.EventsHandlers[eventName], h)
}

func (d eventDispatcher) Dispatch(events []Event) error {
	for _, e := range events {
		handlers, exists := d.EventsHandlers[e.Name()]
		if !exists {
			return ErrNoneHandler
		}
		for _, h := range handlers {
			if err := h.Handle(e); err != nil {
				return err
			}
		}
	}

	return nil
}
