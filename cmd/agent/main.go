package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type MetricForExport struct {
	name       string
	value      string
	stringType string
}

func (m MetricForExport) String() string {
	return m.name + ":" + m.value + ":" + m.stringType
}

type Poller struct {
	PollCount   uint32
	RandomValue float64
	Metrics     []MetricForExport
}

func (p *Poller) fetchMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	metrics := reflect.ValueOf(ms)

	mfe := make([]MetricForExport, 0, metrics.NumField())

	for i := 0; i < metrics.NumField(); i++ {
		name := metrics.Type().Field(i).Name
		value := metrics.Field(i)
		var strValue, stringType string
		switch value.Kind() {
		case reflect.Uint32, reflect.Uint64:
			strValue = strconv.FormatUint(value.Uint(), 10)
			stringType = "counter"
		case reflect.Float64:
			strValue = strconv.FormatFloat(value.Float(), 'f', -1, 64)
			stringType = "gauge"
		default:
			continue
		}

		mfe = append(mfe, MetricForExport{name, strValue, stringType})
	}

	p.RandomValue = rand.ExpFloat64()
	p.PollCount++
	mfe = append(mfe, MetricForExport{
		name:       "PollCount",
		value:      strconv.Itoa(int(p.PollCount)),
		stringType: "counter",
	})
	mfe = append(mfe, MetricForExport{
		name:       "RandomValue",
		value:      strconv.FormatFloat(p.RandomValue, 'f', -1, 64),
		stringType: "gauge",
	})
	p.Metrics = mfe
}

func (p *Poller) printMetrics() {
	for _, m := range p.Metrics {
		fmt.Println(m)
	}
	fmt.Println(len(p.Metrics))
}

func sendMetric(metric *MetricForExport) error {
	urlString := fmt.Sprintf("http://localhost/update/%s/%s/%s", metric.stringType, metric.name, metric.value)
	resp, err := http.Post(urlString, "Content-Type: text/plain", nil)

	fmt.Printf("Sent metric %s to url %s", metric.name, urlString)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Bad status code")
	}

	return nil
}

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	poller := Poller{}
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	for {
		select {
		case <-tickerPoll.C:
			poller.fetchMetrics()
			//poller.printMetrics()
		case <-tickerReport.C:
			fmt.Println("Sending report")
			for _, metric := range poller.Metrics {
				go sendMetric(&metric)
			}
		}
	}
}
