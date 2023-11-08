package log_agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Admin struct {
	client *kadm.Client
}

func NewAdmin(brokers []string) *Admin {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}
	admin := kadm.NewClient(client)
	return &Admin{client: admin}
}
func (a *Admin) TopicExists(topic string) bool {
	ctx := context.Background()
	topicsMetadata, err := a.client.ListTopics(ctx)
	if err != nil {
		panic(err)
	}
	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true
		}
	}
	return false
}
func (a *Admin) CreateTopic(topic string) {
	ctx := context.Background()
	resp, err := a.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		panic(err)
	}
	for _, ctr := range resp {
		if ctr.Err != nil {
			fmt.Printf("Unable to create topic '%s': %s", ctr.Topic, ctr.Err)
		} else {
			fmt.Printf("Created topic '%s'\n", ctr.Topic)
		}
	}
}
func (a *Admin) Close() {
	a.client.Close()
}

type Producer[T any] struct {
	client *kgo.Client
	topic  string
}

func NewProducer[T any](brokers []string, topic string) *Producer[T] {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}
	return &Producer[T]{client: client, topic: topic}
}
func (p *Producer[T]) SendMessage(m T) {
	ctx := context.Background()
	b, _ := json.Marshal(m)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, nil)
}
func (p *Producer[T]) Write(data []byte) (n int, err error) {
	p.client.Produce(context.Background(), &kgo.Record{Topic: p.topic, Value: data}, nil)
	return len(data), nil
}

func (p *Producer[T]) Close() {
	p.client.Close()
}

type Consumer[T any] struct {
	client *kgo.Client
	topic  string
}

func NewConsumer[T any](brokers []string, topic string, groupID string) *Consumer[T] {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		panic(err)
	}
	return &Consumer[T]{client: client, topic: topic}
}
func (c *Consumer[T]) PrintMessage(ctx context.Context) {
	for {
		fetches := c.client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var msg T
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}
			fmt.Println(msg)
		}
	}
}

func (c *Consumer[T]) Close() {
	c.client.Close()
}
