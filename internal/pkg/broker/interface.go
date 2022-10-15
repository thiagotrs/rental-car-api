package broker

type Subscriber interface {
	Subscribe(topic string) <-chan interface{}
}

type Publisher interface {
	Publish(topic string, data interface{})
}

type Broker interface {
	Subscriber
	Publisher
}

type Consumer interface {
	Consume(data interface{})
}

type ConsumerFunc func(data interface{})

func (c ConsumerFunc) Consume(data interface{}) {
	c(data)
}

func Consume(channel <-chan interface{}, consumer Consumer) {
	for data := range channel {
		go func(data interface{}) {
			consumer.Consume(data)
		}(data)
	}
}
