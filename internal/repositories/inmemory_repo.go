package repositories

import (
	"context"
	"fmt"
	"github.com/RomanAVolodin/metrix-go/internal/entities"
	"strconv"
)

type Gauge float64
type Counter int64

type InMemoryRepository struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func (r *InMemoryRepository) Save(ctx context.Context, metric *entities.Metric) error {
	switch metric.MetricType {
	case entities.MetricCounter:
		val, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			return err
		}
		r.Gauges[metric.Name] = Gauge(val)
	case entities.MetricGauge:
		val, err := strconv.ParseInt(metric.Value, 10, 64)
		if err != nil {
			return err
		}
		r.Counters[metric.Name] += Counter(val)
	}
	fmt.Println(r.Gauges)
	fmt.Println(r.Counters)
	return nil
}
