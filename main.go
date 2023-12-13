package main

import "hadoop_metric_exporter/server"

func main() {
	exporter := server.InitlExporterServer()
	exporter.ExposeMetrics()
	exporter.Listen()
}
