package eventbus

type SubscriberManager struct {
	subscriberHandlers map[string][]SubscriberHandler
}

func NewSubscriberManager() *SubscriberManager {
	s := &SubscriberManager{subscriberHandlers: make(map[string][]SubscriberHandler)}
	return s
}

func (sm *SubscriberManager) AddSubscription(event string, sh SubscriberHandler) {
	if !sm.HasSubscription(event) {
		sm.subscriberHandlers[event] = []SubscriberHandler{sh}
	}
	//For now the same event cant have more than one handler
}

func (sm *SubscriberManager) HasSubscription(event string) bool {
	_, ok := sm.subscriberHandlers[event]
	return ok
}

func (sm *SubscriberManager) GetSubscriptionHandlers(event string) []SubscriberHandler {
	return sm.subscriberHandlers[event]
}
