package cart

import "context"

const EventsTopic = "cart.events.topic"

type EventDispatcher interface {
	Dispatch(ctx context.Context, topic, key string, message Event) error
}

type Event interface {
	GetType() string
	GetOccurredAt() string
}

type EventList []Event

func (e *EventList) Record(event Event) {
	*e = append(*e, event)
}

func (e *EventList) Clear() {
	*e = []Event{}
}
