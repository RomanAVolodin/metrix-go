package main

import (
	"github.com/RomanAVolodin/metrix-go/internal/config"
	"github.com/RomanAVolodin/metrix-go/internal/entities"
	"time"
)

func main() {
	poller := entities.Poller{}
	tickerPoll := time.NewTicker(config.PollInterval)
	tickerReport := time.NewTicker(config.ReportInterval)
	for {
		select {
		case <-tickerPoll.C:
			poller.FetchMetrics()
		case <-tickerReport.C:
			for _, metric := range poller.Metrics {
				go metric.SendToServer()
			}
		}
	}
}
