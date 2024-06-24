package handlers

import (
	"fmt"
	"github.com/RomanAVolodin/metrix-go/internal/repositories"
	"net/http"
	"strings"
)

type MetricsHandler struct {
	Repository repositories.MetricsRepository
}

func (h MetricsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusBadRequest)
		return
	}

	path := strings.Split(request.URL.Path, "/")
	if len(path) != 3 {
		http.Error(writer, "Check path string", http.StatusBadRequest)
		return
	}
	metric_type, metric_name, metric_value := strings.ToLower(path[0]), strings.ToLower(path[1]), strings.ToLower(path[2])
	fmt.Println(metric_type, metric_name, metric_value)

	//metric := entities.Metric{
	//	Name:       metric_name,
	//	Value:      metric_value,
	//	MetricType: metric_type,
	//}
	//h.Repository.Save(context.Background(), &metric)

	writer.Header().Set("content-type", "Content-Type: text/plain")
	writer.WriteHeader(http.StatusOK)
}
