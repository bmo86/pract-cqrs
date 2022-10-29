package main

import (
	"cqrs/events"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NastAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	hub := NewHub()
	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NastAddress))

	if err != nil {
		log.Fatalf("%v", err)
	}

	err = n.OnCreatedFeed(func(cfm events.CreatedFeedMessage) {
		hub.BroadCast(NewCreatedFeedMsg(cfm.Id, cfm.Title, cfm.Description, cfm.CreatedAt), nil)
	})

	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()

	go hub.Run()

	http.HandleFunc("/ws", hub.HandlerWs)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
