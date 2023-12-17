package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliverChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("Hello Kafka GO!", "teste", producer, []byte("transferencia"), deliverChan)
	go DeliveryReport(deliverChan) // outra thread de forma async por causa do 'go' na frente
	// e := <-deliverChan
	// msg := e.(*kafka.Message)
	// if msg.TopicPartition.Error != nil {
	// 	fmt.Println("Erro ao enviar")
	// } else {
	// 	fmt.Println("Mensagem enviada:", msg.TopicPartition)
	// }

	producer.Flush(2000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{

		"bootstrap.servers":   "kafka-kafka-1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}

	p, err := kafka.NewProducer(configMap)

	if err != nil {
		log.Println(err.Error())
	}

	return p
}

func Publish(message string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	msg := &kafka.Message{
		Value:          []byte(message),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}

	err := producer.Produce(msg, deliveryChan)

	if err != nil {
		return err
	}

	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar:", ev.TopicPartition.Error)
			} else {
				fmt.Println("Mensagem enviada:", ev.TopicPartition)
			}
		}
	}
}
