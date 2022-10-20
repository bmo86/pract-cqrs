package main

import (
	"cqrs/database"
	"cqrs/events"
	"cqrs/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NastAddress      string `envconfig:"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", handlerCreatedFeed).Methods(http.MethodPost)
	return
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	repo, err := database.NewPostgresRepository(addr)

	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepo(repo) //conn for database

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NastAddress))

	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n) //conn nats

	defer events.Close()

	r := newRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
