package main

import (
	cs "NatsStream_Service/internal/cashe"
	"NatsStream_Service/internal/config"
	"NatsStream_Service/internal/model"
	postgres "NatsStream_Service/internal/storage"
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

func main() {
	cfg := config.MustLoad("SUBSCRIBER") // загрузка конфига subscribera

	// подключение к базе данных
	storagePath := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.User,
		cfg.DataBase.Password, cfg.DataBase.Dbname, cfg.DataBase.Sslmode,
	)
	storage, err := postgres.NewDB(storagePath)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	// создание кэша и его загрузка из бд
	cashe := cs.NewCashe()
	err = storage.UploadCashe(&cashe)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	dataRecieved := *new(model.Order_client)
	fmt.Println("First dataRec", dataRecieved)

	// подключение к Nats-Streaming
	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID)
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &dataRecieved)
		if err != nil {
			log.Println(err)
		}
		// добавление в кэш
		err = cashe.InsertToCashe(dataRecieved) // два одинковых id
		if err != nil {
			log.Println(err)
		}
		// добавление в бд
		err = storage.SaveOrder(dataRecieved)
		if err != nil {
			log.Println(err)
		}
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
	server := &http.Server{Addr: cfg.HTTPServer.Address}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	sub.Unsubscribe()
	sc.Close()
}
