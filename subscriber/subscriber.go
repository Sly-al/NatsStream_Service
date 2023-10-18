package main

import (
	st "NatsStream_Service/internal/cashe"
	"NatsStream_Service/internal/model"
	postgres "NatsStream_Service/internal/storage"
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

const (
	clusterID = "test-cluster"
	clientID  = "Person2"
)

func main() {
	i := 0
	storage, err := postgres.NewDB("host=localhost port=5432 user=user password=userwb dbname=db_wb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	cashe := st.NewCashe()
	err = storage.UploadCashe(&cashe)
	if err != nil {
		log.Fatal(err)
	}

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
		err = cashe.InsertToCashe(dataRecieved) // два одинковых id
		if err != nil {
			log.Println(err)
		}
		err = storage.SaveOrder(dataRecieved)
		if err != nil {
			log.Println(err)
		}
		i++
	}, stan.StartWithLastReceived())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Println("GET METHOD")
			tmpl, err := template.ParseFiles("template/ui.html")
			if err != nil {
				fmt.Println(err)
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				fmt.Println(err)
			}
		case "POST":
			fmt.Println("POST")
			if order, err := cashe.GetFromCashe(r.PostFormValue("order_uid")); err == nil {
				jsonToSend, err := json.MarshalIndent(order, "", " ")
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Sending order with id: %s\n", r.PostFormValue("order_uid"))
					fmt.Fprintf(w, string(jsonToSend))
				}
			} else {
				fmt.Fprint(w, "No such order_uid", err)
			}
			//fmt.Println("POST METHOD")
			//fmt.Fprintf(w, "Hello world!")
		}
	})
	server := &http.Server{Addr: ":8080"}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	sub.Unsubscribe()
	sc.Close()
}
