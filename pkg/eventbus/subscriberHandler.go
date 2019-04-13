package eventbus

type SubscriberHandler interface {
	Handle(event Event)
}
