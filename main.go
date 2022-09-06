package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// variables to start the server
var listenAddress = flag.String("web.listen-address", ":9888", "Address to listen on for web interface.")
var metricPath = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")

func serverMetrics(listenAddress, metricsPath string) error {
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>Custom Resource Exporter</title></head>
			<body>
			<h1>Possible Endpoints</h1>
			<p><a href='` + metricsPath + `'>Prometheus Metrics</a></p>
			</body>
			</html>
		`))
	})
	return http.ListenAndServe(listenAddress, nil)
}

// start the service
func main() {
	log.Fatal(serverMetrics(*listenAddress, *metricPath))
}
