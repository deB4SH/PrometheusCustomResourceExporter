package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	config "github.com/deb4sh/PrometheusCustomResourceExporter/config"
	k8api "github.com/deb4sh/PrometheusCustomResourceExporter/kubernetes"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// variables to start the server
var listenAddress = flag.String("web.listen-address", ":9888", "Address to listen on for web interface.")
var metricPath = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
var configPath = flag.String("config.path", "example.config.yaml", "Path under which the config is located.")

// func serverMetrics(listenAddress, metricsPath string) error {
// 	http.Handle(metricsPath, promhttp.Handler())
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte(`
// 			<html>
// 			<head><title>Custom Resource Exporter</title></head>
// 			<body>
// 			<h1>Possible Endpoints</h1>
// 			<p><a href='` + metricsPath + `'>Prometheus Metrics</a></p>
// 			</body>
// 			</html>
// 		`))
// 	})
// 	return http.ListenAndServe(listenAddress, nil)
// }

// start the service
func main() {
	configValidationErr := config.ValidateConfigFilePath(*configPath)
	if configValidationErr != nil {
		fmt.Println(configValidationErr.Error())
		os.Exit(-1)
	}
	crdConfig, crdConfigErr := config.NewConfig(*configPath)
	if crdConfigErr != nil {
		fmt.Println(crdConfigErr.Error())
		os.Exit(-1)
	}
	//currently for debug purpses
	connection, err := k8api.BuildKubernetesConnection()
	if err != nil {
		fmt.Println("Could not build Kubernetes Connection", err)
		os.Exit(-1)
	}
	//read all configured custom resources
	for _, element := range crdConfig.CustomResourceEntryArray { //not needing the index
		tmpCrd := k8api.NewCustomResourceFromConfigEntry(element)
		cr, _ := k8api.ParseCR(tmpCrd, *connection)
		//unmarschal data of cr
		dynamic := make(map[string]interface{})
		json.Unmarshal([]byte(cr.Data), &dynamic)
		//interprete spec.data  and generate prometheus values from it
		specdata := dynamic["spec"].(map[string]interface{})
		for key, data := range specdata {
			fmt.Println(key)
			fmt.Println(data)
		}
	}
}
