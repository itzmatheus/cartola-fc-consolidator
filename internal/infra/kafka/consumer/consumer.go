package consumer

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/kafka/factory"
	"github.com/itzmatheus/cartola-fc-consolidator/pkg/uow"
)

func Consume(topics []string, servers string, msgChan chan *kafka.Message) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          "gostats",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	kafkaConsumer.SubscribeTopics(topics, nil)
	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		} else {
			log.Printf("Error kafka subscription: %s", err.Error())
		}
	}
}

func ProcessEvents(ctx context.Context, msgChan chan *kafka.Message, uow uow.UowInterface) {
	for msg := range msgChan {
		log.Println("Received Message:", string(msg.Value), "on topic", *msg.TopicPartition.Topic)
		strategy := factory.CreateProcessMessageStrategy(*msg.TopicPartition.Topic)
		err := strategy.Process(ctx, msg, uow)
		if err != nil {
			log.Println("Error Process Message:", string(msg.Value), "on topic", *msg.TopicPartition.Topic, "Erro", err)
		}
	}
}
