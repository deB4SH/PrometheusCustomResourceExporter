package prometheus

import (
	"fmt"
	"log"
	"net/http"
)

type MetricsManager struct {
	Path   string
	Port   int
	Server *http.ServeMux
}

func NewMetricsManager(path string, port int) *MetricsManager {
	p := new(MetricsManager)
	p.Path = path
	p.Port = port

	mux := http.NewServeMux()
	p.Server = mux

	return p
}

func (this *MetricsManager) AddMetricPath(path string, handler http.Handler) {
	this.Server.Handle(path, handler)
}

func (this *MetricsManager) Serve() {
	// s := &http.Server{
	// 	Addr:           fmt.Sprintf(":%d", this.Port),
	// 	Handler:        nil,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// this.Server = s //hand over ref to struct

	//start goroutine so server is running in a seperate thread
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", this.Port), this.Server))
	}()
}
