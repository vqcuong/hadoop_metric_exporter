package server

import "gopkg.in/alecthomas/kingpin.v2"

type ExporterArgs struct {
	Config             *string
	ClusterName        *string
	NameNodeJmx        *string
	DataNodeJmx        *string
	JournalNodeJmx     *string
	ResourceManagerJmx *string
	NodeManagerJmx     *string
	JobHistoryJmx      *string
	HMasterJmx         *string
	HRegionServerJmx   *string
	HiveServer2Jmx     *string
	HiveLLAPJmx        *string
	AutoDiscovery      *bool
	DiscoveryWhileList *string
	Address            *string
	Port               *int
	MetricPath         *string
	LogLevel           *string
}

func ParseExporterArgs() ExporterArgs {
	args := ExporterArgs{
		Config: kingpin.Flag(
			"config",
			"Exporter config file. Default: /exporter/config.yaml",
		).Default("/exporter/config.yaml").Short('c').OverrideDefaultFromEnvar("HADOOP_EXPORTER_CONFIG").String(),
		ClusterName: kingpin.Flag(
			"clusterName",
			"Hadoop cluster labels. Default: hadoop_cluster",
		).Default("hadoop_cluster").Short('n').OverrideDefaultFromEnvar("HADOOP_EXPORTER_CLUSTER_NAME").String(),
		NameNodeJmx: kingpin.Flag(
			"nameNodeJmx",
			"List of HDFS namenode JMX url seperated by comma (,). Example: http://localhost:9870/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_NAMENODE_JMX").String(),
		DataNodeJmx: kingpin.Flag(
			"dataNodeJmx",
			"List of HDFS datanode JMX url seperated by comma (,). Example: http://localhost:9864/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_DATANODE_JMX").String(),
		JournalNodeJmx: kingpin.Flag(
			"journalNodeJmx",
			"List of HDFS journalnode JMX url seperated by comma (,). Example: http://localhost:8480/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_JOURNALNODE_JMX").String(),
		ResourceManagerJmx: kingpin.Flag(
			"resourceManagerJmx",
			"List of YARN resourcemanager JMX url seperated by comma (,). Example: http://localhost:8088/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_RESOURCEMANAGER_JMX").String(),
		NodeManagerJmx: kingpin.Flag(
			"nodeManagerJmx",
			"List of YARN nodemanager JMX url seperated by comma (,). Example: http://localhost:8042/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_NODEMANAGER_JMX").String(),
		JobHistoryJmx: kingpin.Flag(
			"jobHistoryJmx",
			"List of Mapreduce jobhistory JMX url seperated by comma (,). Example: http://localhost:19888/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_JOBHISTORY_JMX").String(),
		HMasterJmx: kingpin.Flag(
			"hMasterJmx",
			"List of HBase master JMX url seperated by comma (,). Example: http://localhost:16010/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_HMASTER_JMX").String(),
		HRegionServerJmx: kingpin.Flag(
			"hRegionServerJmx",
			"List of HBase regionserver JMX url seperated by comma (,). Example: http://localhost:16030/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_HREGIONSERVER_JMX").String(),
		HiveServer2Jmx: kingpin.Flag(
			"hiveServer2Jmx",
			"List of HiveServer2 JMX url seperated by comma (,). Example: http://localhost:10002/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_HIVESERVER2_JMX").String(),
		HiveLLAPJmx: kingpin.Flag(
			"hiveLLAPJmx",
			"List of Hive LLAP JMX url seperated by comma (,). Example: http://localhost:15002/jmx",
		).OverrideDefaultFromEnvar("HADOOP_EXPORTER_HIVELLAP_JMX").String(),
		AutoDiscovery: kingpin.Flag(
			"autoDiscovery",
			"Enable auto discovery if set true else false. Default: false",
		).Default("false").Short('a').OverrideDefaultFromEnvar("HADOOP_EXPORTER_AUTO_DISCOVERY").Bool(),
		DiscoveryWhileList: kingpin.Flag(
			"discoveryWhileList",
			"List of shortnames of services (namenode: nn, datanode: dn, ...) that should be enable to auto discovery",
		).Short('w').OverrideDefaultFromEnvar("HADOOP_EXPORTER_DISCOVERY_WHITELIST").String(),
		Address: kingpin.Flag(
			"address",
			"Polling server on this address (hostname or ip). Default: 0.0.0.0",
		).Default("0.0.0.0").Short('h').OverrideDefaultFromEnvar("HADOOP_EXPORTER_ADDRESS").String(),
		Port: kingpin.Flag(
			"port",
			"Port to listen on. Default: 9123",
		).Default("9123").Short('p').OverrideDefaultFromEnvar("HADOOP_EXPORTER_PORT").Int(),
		MetricPath: kingpin.Flag(
			"metricPath",
			"Path under which to expose metrics. Default: /metrics",
		).Default("/metrics").Short('d').OverrideDefaultFromEnvar("HADOOP_EXPORTER_METRIC_PATH").String(),
		LogLevel: kingpin.Flag(
			"logLevel",
			"Log level, include: all, debug, info, warn, error. Default: info",
		).Default("info").Short('l').OverrideDefaultFromEnvar("HADOOP_EXPORTER_LOG_LEVEL").String(),
	}
	kingpin.Parse()
	return args
}
