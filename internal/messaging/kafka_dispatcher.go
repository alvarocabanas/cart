package messaging

import (
	"context"
	"encoding/json"

	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"

	"github.com/Shopify/sarama"
	cart "github.com/alvarocabanas/cart/internal"
)

type KafkaEventDispatcher struct {
	producer sarama.SyncProducer
	client   sarama.Client
}

func NewKafkaProducer(kafkaBrokers KafkaBrokers) (*KafkaEventDispatcher, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Return.Errors = true
	kafkaCfg.Version = sarama.V1_0_0_0
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaCfg.Producer.Return.Successes = true

	kafkaClient, err := sarama.NewClient(kafkaBrokers, kafkaCfg)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		return nil, err
	}

	return &KafkaEventDispatcher{producer: producer, client: kafkaClient}, nil
}

func (p *KafkaEventDispatcher) Dispatch(ctx context.Context, topic, key string, event cart.Event) error {
	_, span := trace.StartSpan(ctx, "dispatch_add_item_event")
	defer span.End()

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
		Key:   sarama.ByteEncoder(key),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("traceSpan"),
				Value: propagation.Binary(span.SpanContext()),
			},
		},
	}
	_, _, err = p.producer.SendMessage(producerMessage)

	return err
}
