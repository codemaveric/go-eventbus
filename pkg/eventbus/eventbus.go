package eventbus

type EventBus interface {
	Publish(event Event)
	Subscribe(event Event, subscriberHandler SubscriberHandler)
}
