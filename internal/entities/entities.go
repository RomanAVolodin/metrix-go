package entities

import (
	"errors"
	"fmt"
	"github.com/RomanAVolodin/metrix-go/internal/config"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

type MetricForExport struct {
	Name       string
	Value      string
	StringType string
}

func (m MetricForExport) String() string {
	return m.Name + ":" + m.Value + ":" + m.StringType
}

func (m *MetricForExport) SendToServer() error {
	urlString := fmt.Sprintf("%s/update/%s/%s/%s", config.ServerHost, m.StringType, m.Name, m.Value)
	resp, err := http.Post(urlString, "Content-Type: text/plain", nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

type Poller struct {
	PollCount   uint32
	RandomValue float64
	Metrics     []MetricForExport
}

func (p *Poller) FetchMetrics() {
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
		Name:       "PollCount",
		Value:      strconv.Itoa(int(p.PollCount)),
		StringType: "counter",
	})
	mfe = append(mfe, MetricForExport{
		Name:       "RandomValue",
		Value:      strconv.FormatFloat(p.RandomValue, 'f', -1, 64),
		StringType: "gauge",
	})
	p.Metrics = mfe
}

func (p *Poller) PrintMetrics() {
	for _, m := range p.Metrics {
		fmt.Println(m)
	}
	fmt.Println(len(p.Metrics))
}
