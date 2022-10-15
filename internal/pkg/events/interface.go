package events

type Event interface {
	Name() string
}

type EventHandlerFunc func(e Event) error

func (h EventHandlerFunc) Handle(e Event) error {
	return h(e)
}

type EventHandler interface {
	Handle(e Event) error
}

type Dispatcher interface {
	Register(h EventHandler, eventName string)
	Dispatch(events []Event) error
}
