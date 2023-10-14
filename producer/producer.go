package main

import (
	"NatsStream_Service/internal/model"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
)

func main() {
	clusterID := "test-cluster"
	clientID := "Person1"

	order_uids := [3]string{
		"first",
		"second",
		"third",
	}

	json_byte_read, err := os.ReadFile("internal/model/model.json")
	if err != nil {
		log.Fatalf("Unamle to read due to %s\n", err)
	}

	var order model.Order_client

	err = json.Unmarshal(json_byte_read, &order)
	if err != nil {
		log.Fatalf("Unable to unmarshal JSON due to %s", err)
	}

	sc, err := stan.Connect(clusterID, clientID)
	defer sc.Close()

	if err != nil {
		log.Print(err)
		return
	}

	for _, now := range order_uids {
		order.Order_uid = now
		json_byte_send, err := json.MarshalIndent(order, "", " ")
		if err != nil {
			log.Fatalf("Unable to marshal JSON due to %s", err)
		}
		sc.Publish("foo", json_byte_send)
		time.Sleep(time.Second * 3)
	}

}
