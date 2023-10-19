package main

import (
	"NatsStream_Service/internal/config"
	"NatsStream_Service/internal/model"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

func main() {
	cfg := config.MustLoad("PRODUCER") // загрузка конфига для producerа
	var order model.Order_client       // структура заказа

	// подключение к Nats-Streaming
	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID)
	if err != nil {
		log.Fatalf("Unable to connect %s", err)
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 5)

		// Создание новых заказов
		err = gofakeit.Struct(&order)
		if err != nil {
			log.Printf("Unable to generate json due to %s", err)
		}

		// Order_client -> []byte
		jsonToSend, err := json.MarshalIndent(order, "", " ")
		if err != nil {
			log.Printf("Unable to marshal JSON due to %s", err)
		}

		// Публикация в канал
		err = sc.Publish("foo", jsonToSend)
		if err != nil {
			log.Printf("Json wasn't published due to: %s", err)
		} else {
			log.Printf("Send successfully: %s", order.Order_uid)
		}
	}

	sc.Close()
}
