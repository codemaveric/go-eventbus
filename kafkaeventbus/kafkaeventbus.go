package kafkaeventbus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/codemaveric/go-eventbus/eventbus"

	"github.com/segmentio/kafka-go"
)

type KafkaEventBus struct {
	subscriberManager *eventbus.SubscriberManager
	config            *KafkaConfig
}

func NewKafkaEventBus(subscriberManager *eventbus.SubscriberManager, config *KafkaConfig) *KafkaEventBus {
	k := &KafkaEventBus{
		subscriberManager: subscriberManager,
		config:            config,
	}
	return k
}

func (kb *KafkaEventBus) Publish(event eventbus.Event) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kb.config.Server},
		Topic:    event.GetName(),
		Balancer: &kafka.LeastBytes{},
	})
	payload, _ := json.Marshal(event)
	w.WriteMessages(context.Background(), kafka.Message{Value: payload})
	w.Close()
}

func (kb *KafkaEventBus) Subscribe(event eventbus.Event, handler eventbus.SubscriberHandler) {

	hasSub := kb.subscriberManager.HasSubscription(event.GetName())
	//Add Subscription to manager, if subscription exist just register the handler if its a new handler
	kb.subscriberManager.AddSubscription(event.GetName(), handler)
	if !hasSub {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{kb.config.Server},
			GroupID:  kb.config.ConsumerConfig.Group,
			Topic:    event.GetName(),
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
		go kb.process(event, r)
	}

}

func (kb *KafkaEventBus) process(event eventbus.Event, r *kafka.Reader) {
	ctx := context.Background()
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Fatal(err)
			continue
		}
		handlers := kb.subscriberManager.GetSubscriptionHandlers(event.GetName())
		for _, handler := range handlers {
			handler.Handle(event) //Handle
		}
		r.CommitMessages(ctx, msg)
	}
	r.Close()
}
