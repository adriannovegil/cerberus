package prometheus

import (
	"net/http"

	cPrometheus "github.com/prometheus/client_golang/prometheus"
	cPromHttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config data structure
type Config struct {
	Enable bool   `yaml:"enable"`
	Port   int    `yaml:"port"`
	Path   string `yaml:"path"`
}

// StartPrometheusServer start the Prometheus service endpoint
func StartPrometheusServer() cPrometheus.Registerer {
	// Prometheus registry to expose metrics.
	registerer := cPrometheus.NewRegistry()

	go func() {
		http.ListenAndServe(":8081", cPromHttp.HandlerFor(registerer, cPromHttp.HandlerOpts{}))
	}()

	return registerer
}
