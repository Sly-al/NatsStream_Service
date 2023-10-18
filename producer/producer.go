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
	cfg := config.MustLoad("PRODUCER")
	var order model.Order_client // инициализая структуры заказа
	
	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID) // подключение к Nats-Streaming
	defer sc.Close()

	if err != nil {
		log.Print(err)
		return
	}

	for i := 0; i < 3; i++ {
		gofakeit.Struct(&order)                               // Создание новых заказов
		jsonToSend, err := json.MarshalIndent(order, "", " ") // Order_client -> []byte
		if err != nil {
			log.Fatalf("Unable to marshal JSON due to %s", err)
		}
		sc.Publish("foo", jsonToSend) // Публикация в канал
		time.Sleep(time.Second * 5)
	}

}
