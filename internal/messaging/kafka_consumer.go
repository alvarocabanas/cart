package messaging

import (
	"context"
	"errors"
	"log"

	"github.com/Shopify/sarama"
)

type (
	KafkaBrokers         []string
	KafkaConsumerGroupID string
	KafkaTopics          []string
	HandleFunc           func(context.Context, []byte) error
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	consumer      consumer
	client        sarama.Client
}

func NewKafkaConsumer(
	kafkaBrokers KafkaBrokers,
	consumerGroupID KafkaConsumerGroupID,
	topics KafkaTopics,
	handleFunc HandleFunc,
) (*KafkaConsumer, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Return.Errors = true
	kafkaCfg.Version = sarama.V1_0_0_0
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewClient(kafkaBrokers, kafkaCfg)
	if err != nil {
		return nil, err
	}

	consumerGroup, err := sarama.NewConsumerGroupFromClient(string(consumerGroupID), client)
	if err != nil {
		return nil, err
	}

	consumer := consumer{
		topics:     topics,
		handleFunc: handleFunc,
	}

	return &KafkaConsumer{
		consumer:      consumer,
		consumerGroup: consumerGroup,
		client:        client,
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	defer func() {
		_ = c.consumerGroup.Close()
	}()

	for {
		select {
		case err := <-c.consumerGroup.Errors():
			return err
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				return ctx.Err()
			}

			return nil
		default:
			if err := c.consumerGroup.Consume(ctx, c.consumer.topics, &c.consumer); err != nil {
				return err
			}
		}
	}
}

type consumer struct {
	topics     []string
	handleFunc func(context.Context, []byte) error
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		select {
		case <-session.Context().Done():
			return nil
		default:
			err := c.handleFunc(session.Context(), message.Value)
			log.Print(err)

			session.MarkMessage(message, "")
		}
	}

	return nil
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
