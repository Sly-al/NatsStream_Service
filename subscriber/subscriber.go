package main

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"os"
	"time"
)

func main() {
	var json_recieved []byte
	clusterID := "test-cluster"
	clientID := "Person2"

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print("Connected")

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		json_recieved = m.Data
	}, stan.StartWithLastReceived())

	time.Sleep(time.Second)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("%d%s", i, ".json")
		err = os.WriteFile(name, json_recieved, 0644)
		time.Sleep(time.Second * 2)
	}

	// Unsubscribe
	sub.Unsubscribe()

	// Close connection
	sc.Close()

}
