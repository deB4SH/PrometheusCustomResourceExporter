package prometheus

import (
	"fmt"
	"net/http"
)

type MetricsManager struct {
	Path string
	Port int
	Http http.Server
}

func New(path string, port int) *MetricsManager {
	p := new(MetricsManager)
	p.Path = path
	p.Port = port
	p.Http = http.DefaultServeMux
	return p
}

func addMetricPath(m MetricsManager, path string, handler http.Handler) {
	m.Http.Handle()
}

func Serve(m MetricsManager) {
	http.ListenAndServe(fmt.Sprintf(":%d", m.Port), nil)
}
