package main

import (
	"log"
	"net/http"
	"workshop/internal/api/jokes"
	"workshop/internal/config"

	"github.com/go-chi/chi"
	"github.com/ilyakaznacheev/cleanenv"

	"workshop/internal/handler"
)

func main() {
	cfg := config.Server{}
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	jokeClient := jokes.NewJobClient(cfg.JokeURL)

	h := handler.NewHandler(jokeClient)
	r := chi.NewRouter()
	r.Get("/hello", h.Hello)
	path := cfg.Host+":"+cfg.Port

	log.Printf("starting server at %s", path)
	err = http.ListenAndServe(path, r)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("shutting server down")
}
