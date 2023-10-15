package main

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"os"
	"strconv"
)

const (
	clusterID = "test-cluster"
	clientID  = "Person2"
)

func main() {
	var (
		i    int
		name string
	)

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print("Connected")

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		name = strconv.Itoa(i+1) + ".json"
		err = os.WriteFile(name, m.Data, 0644)
		i++
	}, stan.StartWithLastReceived())

	for i < 3 {
		_ = 0
	}
	sub.Unsubscribe()
	sc.Close()
}
