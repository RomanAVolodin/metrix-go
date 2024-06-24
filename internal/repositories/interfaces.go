package repositories

import (
	"context"
	"github.com/RomanAVolodin/metrix-go/internal/entities"
)

type MetricsRepository interface {
	Save(ctx context.Context, metric *entities.Metric) error
}
