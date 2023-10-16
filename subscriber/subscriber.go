package main

import (
	st "NatsStream_Service/internal/cashe"
	"NatsStream_Service/internal/model"
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"log"
)

const (
	clusterID = "test-cluster"
	clientID  = "Person2"
)

func main() {
	i := 0

	cashe := st.NewCashe()
	dataRecieved := *new(model.Order_client)

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		fmt.Print(err)
		return
	}

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &dataRecieved)
		if err != nil {
			log.Println(err)
		}
		err = cashe.InsertToCashe(dataRecieved)
		if err != nil {
			log.Println(err)
		}
		i++
	}, stan.StartWithLastReceived())

	for i < 3 {
		_ = 0
	}

	for key, val := range cashe.Data {
		fmt.Println(key)
		fmt.Println(val)
	}
	sub.Unsubscribe()
	sc.Close()
}
