# Hadoop Prometheus Exporter

An hadoop metrics exporter for common hadoop components. Currently, I've just implemented for HDFS NameNode, HDFS DataNode, HDFS JournalNode, YARN ResourceManager, YARN NodeManager. This is a golang version, you may take another version using python [here](https://github.com/vqcuong/hadoop_exporter).

## How it works

- Consume metrics from JMX http, convert and export hadoop metrics via HTTP for Prometheus consumption.
- Underlyring, I used regex template to parse and map config name as well as label before exporting it via promethues http server. You can see my templates in folder [metrics](./metrics)

## Install go environment

- Follow this [docs](https://github.com/go-nv/goenv/blob/master/INSTALL.md) to install goenv
- Install specified go version

```
# the using version is defined in .go-version file
goenv install --skip-existing
```

## How to run

```
go build && ./hadoop_metric_exporter
```

Help on flags of hadoop_metric_exporter:

```
$ hadoop_metric_exporter --help
usage: hadoop_metric_exporter [<flags>]
Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="/exporter/config.yaml"
                                 Exporter config file. Default: /exporter/config.yaml
  -n, --clusterName="hadoop_cluster"
                                 Hadoop cluster labels. Default: hadoop_cluster
      --nameNodeJmx=NAMENODEJMX  List of HDFS namenode JMX url seperated by comma (,). Example:
                                 http://localhost:9870/jmx
      --dataNodeJmx=DATANODEJMX  List of HDFS datanode JMX url seperated by comma (,). Example:
                                 http://localhost:9864/jmx
      --journalNodeJmx=JOURNALNODEJMX
                                 List of HDFS journalnode JMX url seperated by comma (,). Example:
                                 http://localhost:8480/jmx
      --resourceManagerJmx=RESOURCEMANAGERJMX
                                 List of YARN resourcemanager JMX url seperated by comma (,). Example:
                                 http://localhost:8088/jmx
      --nodeManagerJmx=NODEMANAGERJMX
                                 List of YARN nodemanager JMX url seperated by comma (,). Example:
                                 http://localhost:8042/jmx
      --jobHistoryJmx=JOBHISTORYJMX
                                 List of Mapreduce jobhistory JMX url seperated by comma (,). Example:
                                 http://localhost:19888/jmx
      --hMasterJmx=HMASTERJMX    List of HBase master JMX url seperated by comma (,). Example:
                                 http://localhost:16010/jmx
      --hRegionServerJmx=HREGIONSERVERJMX
                                 List of HBase regionserver JMX url seperated by comma (,). Example:
                                 http://localhost:16030/jmx
      --hiveServer2Jmx=HIVESERVER2JMX
                                 List of HiveServer2 JMX url seperated by comma (,). Example:
                                 http://localhost:10002/jmx
      --hiveLLAPJmx=HIVELLAPJMX  List of Hive LLAP JMX url seperated by comma (,). Example:
                                 http://localhost:15002/jmx
  -a, --autoDiscovery            Enable auto discovery if set true else false. Default: false
  -w, --discoveryWhileList=DISCOVERYWHILELIST
                                 List of shortnames of services (namenode: nn, datanode: dn, ...) that should be
                                 enable to auto discovery
  -h, --address="0.0.0.0"        Enable auto discovery if set true else false. Default: 0.0.0.0
  -p, --port=9123                Listen to this port. Default: 9123
  -d, --metricPath="/metrics"    Path under which to expose metrics. Default: /metrics
  -l, --logLevel="info"          Log level, include: all, debug, info, warn, error. Default: info
```

You can use config file (yaml format) to replace commandline args. Example of config.yaml:

```
# exporter server config
server:
  address: 127.0.0.1 # address to run exporter
  port: 9123 # port to listen
  metricPath: /metrics # metric path to expose
  logLevel: info # logging level included: all, info, warn, debug, error

# list of jmx service to scape metrics
jmx:
  - cluster: hadoop_prod
    services:
      nameNode:
        - http://nn1:9870/jmx
      dataNode:
        - http://dn1:9864/jmx
        - http://dn2:9864/jmx
        - http://dn3:9864/jmx
      resourceManager:
        - http://rm1:8088/jmx
      nodeManager:
        - http://nm1:8042/jmx
        - http://nm2:8042/jmx
        - http://nm3:8042/jmx
      hiveServer2:
        - http://hs2:10002/jmx
      hMaster:
        - http://hmaster1:16010/jmx
        - http://hmaster2:16010/jmx
        - http://hmaster3:16010/jmx
      hRegionserver:
        - http://hregionserver1:16030/jmx
        - http://hregionserver2:16030/jmx
        - http://hregionserver3:16030/jmx

  - cluster: hadoop_dev
    services:
      nameNode:
        - http://dev:9870/jmx
      dataNode:
        - http://dev:9864/jmx
      resourceManager:
        - http://dev:8088/jmx
      nodeManager:
        - http://dev:8042/jmx
```

Tested on Apache Hadoop 2.7.3, 3.3.0, 3.3.1, 3.3.2

## Grafana Monitoring

There are [HDFS](./dashboards/hdfs.json) and [YARN](./dashboards/yarn.json) dashboard definition prepared by me. You can import it directly on grafana.

## Docker deployment

Run container:

```
docker run -d \
  --name hadoop-metric-exporter \
  vqcuong96/hadoop_metric_exporter \
  -nn http://localhost:9870/jmx \
  -rm http://localhost:8088/jmx
```

You can also mount config to docker container:

```
docker run -d \
  --name hadoop-metric-exporter \
  --mount type=bind,source=/path/to/config.yaml,target=/tmp/config.yaml \
  vqcuong96/hadoop_metric_exporter \
  -c /tmp/config.yaml
```

To build your own images, run:

```
./docker/build.sh [your_repo] [your_version_tag]
```

For example:

```
./build.sh mydockerhub/ latest
#your image will look like: mydockerhub/hadoop_metric_exporter:latest
```
