package main

import (
	"cqrs/database"
	"cqrs/events"
	"cqrs/repository"
	"cqrs/search"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NastAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAddress string `envconfig:"ELASTICSERACH_ADDRESS"`
}

func newRouter() (r *mux.Router) {
	r = mux.NewRouter() //instanciar el router
	r.HandleFunc("/feeds", handlerListFeed).Methods(http.MethodGet)
	r.HandleFunc("/search", handlerSearch).Methods(http.MethodGet)
	return
}

func main() {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("%v", err)
	}

	//conn for database
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmodel=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepo(repo)

	//conn elasticSearch
	es, err := search.NewElastic(fmt.Sprintf("http://%s", cfg.ElasticSearchAddress))
	if err != nil {
		log.Fatal(err)
	}

	search.SetSearchRepository(es)

	defer search.Close()

	//conn nats
	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NastAddress))
	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreatedFeed(onCreatedFeed)
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()

	r := newRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
