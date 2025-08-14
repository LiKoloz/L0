package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka" // Убедитесь, что пакет установлен
)

func main() {
	// Конфигурация потребителя
	config := &kafka.ConfigMap{
 "bootstrap.servers": "192.168.18.138:9092,192.168.18.138:9093,192.168.18.138:9094", 
  "group.id":          "myGroup", 
      "auto.offset.reset": "smallest", 
        }
	// Создание потребителя
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	// Подписка на топик
	topic := "test-topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Subscribe failed: %s\n", err)
		os.Exit(1)
	}

	// Обработка сигналов завершения
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Основной цикл обработки сообщений
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			// Таймаут опроса 100мс
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Обработка сообщения
				fmt.Printf("Received message:\n%s\n", string(e.Value))
				fmt.Printf("Headers: %v\n", e.Headers) // Исправлено: e.Headers вместо e.eaders
				fmt.Printf("Partition: %d, Offset: %d\n\n", 
					e.TopicPartition.Partition, e.TopicPartition.Offset)
				
				// Подтверждение обработки (при ручном управлении оффсетами)
				// consumer.CommitMessage(e)
				
			case kafka.Error:
				// Обработка ошибок
				fmt.Printf("Error: %v\n", e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			}
		}
	}
}