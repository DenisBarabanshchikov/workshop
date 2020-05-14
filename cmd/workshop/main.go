package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"workshop/internal/api/jokes"
	"workshop/internal/config"
	"workshop/internal/handler"

	"github.com/go-chi/chi"
	"github.com/ilyakaznacheev/cleanenv"
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
	path := cfg.Host + ":" + cfg.Port

	srv := &http.Server{
		Addr:    path,
		Handler: r,
	}

	//handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		done <- srv.Shutdown(ctx)
	}()

	log.Printf("starting server at %s", path)
	_ = srv.ListenAndServe()

	err = <-done
	log.Printf("shutting server down with %v", err)
}
