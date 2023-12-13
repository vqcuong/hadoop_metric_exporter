package server

import (
	"fmt"
	"hadoop_exporter/collector"
	"hadoop_exporter/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var SERVICE_COMPONENT_MAPPER = map[string]string{
	"namenode":        "hdfs",
	"datanode":        "hdfs",
	"journalnode":     "hdfs",
	"resourcemanager": "yarn",
	"nodemanager":     "yarn",
	"hiveserver2":     "hive",
}

var args = ParseExporterArgs()

var SERVICE_MAPPER = map[string]struct {
	component  string
	argsUrl    string
	defaultUrl string
}{
	"namenode":        {"hdfs", *args.NameNodeJmx, "http://localhost:9870/jmx"},
	"datanode":        {"hdfs", *args.DataNodeJmx, "http://localhost:9864/jmx"},
	"journalnode":     {"hdfs", *args.JournalNodeJmx, "http://localhost:8480/jmx"},
	"resourcemanager": {"yarn", *args.ResourceManagerJmx, "http://localhost:8088/jmx"},
	"nodemanager":     {"yarn", *args.NodeManagerJmx, "http://localhost:8042/jmx"},
	"jobhistory":      {"mapred", *args.JobHistoryJmx, "http://localhost:19888/jmx"},
	"hiveserver2":     {"hive", *args.HiveServer2Jmx, "http://localhost:10002/jmx"},
	"hivellap":        {"hive", *args.HiveLLAPJmx, "http://localhost:15002/jmx"},
	"hmaster":         {"hbase", *args.HMasterJmx, "http://localhost:16010/jmx"},
	"hregionserver":   {"hbase", *args.HRegionServerJmx, "http://localhost:16030/jmx"},
}

type ExporterServer struct {
	config             *ExporterConfig
	address            string
	metricPath         string
	logLevel           string
	collectors         []prometheus.Collector
	discoveryWhitelist []string
	port               int
	autoDiscovery      bool
}

func InitlExporterServer() *ExporterServer {
	server := ExporterServer{
		autoDiscovery:      *args.AutoDiscovery,
		discoveryWhitelist: strings.Split(*args.DiscoveryWhileList, ","),
		address:            *args.Address,
		port:               *args.Port,
		metricPath:         *args.MetricPath,
		logLevel:           *args.LogLevel,
	}

	var config *ExporterConfig
	var err error
	if utils.IsFile(*args.Config) {
		logrus.Infof("Use provided config: %s", *args.Config)
		config, err = ReadExporterConfig(*args.Config)
		if err != nil {
			logrus.Warnf("Something wrong when loading yaml config: %s. Skip it ...", *args.Config)
		}
	}
	if config != nil {
		serverConfig := *config.Server
		server.address = serverConfig.Address
		server.port = serverConfig.Port
		server.metricPath = serverConfig.MetricPath
		server.logLevel = serverConfig.LogLevel
		server.config = config
	}
	utils.Handlelogger(server.logLevel)
	server.handleCollectors()
	return &server
}

func (server *ExporterServer) handleCollectors() {
	if server.config != nil {
		jmx := *server.config.Jmx
		for _, cluster := range jmx {
			newCollectors := buildClusterCollectors(&cluster)
			server.collectors = append(server.collectors, newCollectors...)
		}
		return
	}

	if server.autoDiscovery {
		logrus.Info("Enable service auto discovery mode")
	}

	for service, v := range SERVICE_MAPPER {
		url := v.argsUrl
		if server.autoDiscovery && utils.Contains(server.discoveryWhitelist, service) {
			url = utils.CoalesceString(v.argsUrl, v.defaultUrl)
		}
		if url != "" {
			server.collectors = append(
				server.collectors,
				collector.InitMetricCollector(
					*args.ClusterName,
					v.component,
					service,
					&[]string{url},
				),
			)
		}
	}
}

func buildClusterCollectors(clusterJmx *HadoopClusterJmx) []prometheus.Collector {
	cluster := clusterJmx.Cluster
	collectors := make([]prometheus.Collector, len(clusterJmx.Services))
	i := 0
	for service, urls := range clusterJmx.Services {
		collectors[i] = prometheus.Collector(collector.InitMetricCollector(
			cluster,
			SERVICE_MAPPER[strings.ToLower(service)].component,
			strings.ToLower(service),
			urls,
		))
		i += 1
	}
	return collectors
}

func (server *ExporterServer) ExposeMetrics() {
	prometheus.MustRegister(server.collectors...)
	// http.Handle(server.metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, metricHandler(server.collectors)))
	http.Handle(server.metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, `<html>
			<head><title>Hadoop Exporter</title></head>
			<body>
			<h1>Hadoop Exporter</h1>
			<p><a href="%s">Metrics</a></p>
			</body>
			</html>`, server.metricPath)
	})
}

func (server *ExporterServer) Listen() {
	logrus.Infof("Metrics endpoint - http://%s:%v/metrics", server.address, server.port)
	address := server.address + ":" + strconv.FormatInt(int64(server.port), 10)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		logrus.Fatalf("Something wrong when starting listen on %s: %v", address, err)
	}
}
