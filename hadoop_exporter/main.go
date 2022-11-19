package main

import (
	"hadoop_exporter/server"
)

func main() {
	exporter := server.InitlExporterServer()
	exporter.ExposeMetrics()
	exporter.Listen()
}
