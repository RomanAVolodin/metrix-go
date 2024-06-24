package entities

import (
	"fmt"
	"github.com/RomanAVolodin/metrix-go/internal/config"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

const (
	MetricCounter = "counter"
	MetricGauge   = "gauge"
)

type Metric struct {
	Name       string
	Value      string
	MetricType string
}

func (m Metric) String() string {
	return m.Name + ":" + m.Value + ":" + m.MetricType
}

type Poller struct {
	PollCount   uint32
	RandomValue float64
	Metrics     []*Metric
}

func (p *Poller) FetchMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	metrics := reflect.ValueOf(ms)

	mfe := make([]*Metric, 0, metrics.NumField())

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

		mfe = append(mfe, &Metric{name, strValue, stringType})
	}

	p.RandomValue = rand.ExpFloat64()
	p.PollCount++
	mfe = append(mfe, &Metric{
		Name:       "PollCount",
		Value:      strconv.Itoa(int(p.PollCount)),
		MetricType: MetricCounter,
	})
	mfe = append(mfe, &Metric{
		Name:       "RandomValue",
		Value:      strconv.FormatFloat(p.RandomValue, 'f', -1, 64),
		MetricType: MetricGauge,
	})
	p.Metrics = mfe
}

func (p *Poller) PrintMetrics() {
	for _, m := range p.Metrics {
		fmt.Println(m)
	}
	fmt.Println(len(p.Metrics))
}

func (p *Poller) SendToServer() error {
	for _, m := range p.Metrics {
		go func(m *Metric) {
			urlString := fmt.Sprintf("%s/update/%s/%s/%s", config.ServerHost, m.MetricType, m.Name, m.Value)
			_, _ = http.Post(urlString, "Content-Type: text/plain", nil)

			//if err != nil {
			//	return err
			//}
			//
			//defer resp.Body.Close()
			//
			//if resp.StatusCode != http.StatusOK {
			//	return errors.New("bad status code")
			//}

			fmt.Println(urlString)
		}(m)

	}

	fmt.Println("SENT", len(p.Metrics))
	return nil
}
