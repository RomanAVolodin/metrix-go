package main

import (
	"github.com/RomanAVolodin/metrix-go/internal/config"
	"github.com/RomanAVolodin/metrix-go/internal/handlers"
	"github.com/RomanAVolodin/metrix-go/internal/repositories"
	"log"
	"net/http"
)

func main() {
	handler := handlers.MetricsHandler{Repository: &repositories.InMemoryRepository{
		Gauges:   make(map[string]repositories.Gauge),
		Counters: make(map[string]repositories.Counter),
	}}

	mux := http.NewServeMux()
	mux.Handle(config.UpdateURL, http.StripPrefix(config.UpdateURL, handler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
