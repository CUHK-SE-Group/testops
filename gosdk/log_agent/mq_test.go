package log_agent

import (
	"fmt"
	"testing"
	"time"
)

type testMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func TestNewProducer(t *testing.T) {
	topic := "chat-room"
	brokers := []string{"debian:19092"}
	admin := NewAdmin(brokers)
	defer admin.Close()
	if !admin.TopicExists(topic) {
		admin.CreateTopic(topic)
	}

	producer := NewProducer[testMessage](brokers, topic)
	defer producer.Close()
	consumer := NewConsumer[testMessage](brokers, topic, "test1")
	defer consumer.Close()
	producer.SendMessage(testMessage{
		User:    "Alice",
		Message: "123456",
	})
	fmt.Println("sending message")
	time.Sleep(3 * time.Second)
	//consumer.PrintMessage(context.Background())
}
