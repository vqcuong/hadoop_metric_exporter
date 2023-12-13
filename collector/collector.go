package collector

import (
	"encoding/json"
	"fmt"
	"hadoop_metric_exporter/utils"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var NON_METRIC_NAMES = []string{"name", "modelerType", "Name", "ObjectName"}

type MetricCollector struct {
	rules                map[string][]MetricRulePattern
	commonLabels         map[string]map[string][]string
	firstGetCommonLabels map[string]bool
	metricDescs          *map[string]map[string]*prometheus.Desc
	urls                 *[]string
	cluster              string
	component            string
	service              string
	prefix               string
	isLowerName          bool
	isLowerLabel         bool
}

var EXPORTER_RULES_DIR = utils.CoalesceString(os.Getenv("HADOOP_METRIC_EXPORTER_RULES_DIR"), "./rules")

func InitMetricCollector(cluster string, component string, service string, urls *[]string) MetricCollector {
	commonRuleFile := fmt.Sprintf("%s/%s.yaml", EXPORTER_RULES_DIR, "common")
	serviceRuleFile := fmt.Sprintf("%s/%s.yaml", EXPORTER_RULES_DIR, service)

	commonRules, err1 := ReadServiceRules(commonRuleFile)
	serviceRules, err2 := ReadServiceRules(serviceRuleFile)
	if err1 != nil {
		logrus.Warnf("Unable load service rules: %s", commonRuleFile)
	}
	if err2 != nil {
		logrus.Warnf("Unable load service rules: %s", serviceRuleFile)
	}

	firstGetCommonLabels := make(map[string]bool)
	for _, url := range *urls {
		firstGetCommonLabels[url] = true
	}
	collector := MetricCollector{
		cluster:              cluster,
		component:            component,
		service:              service,
		prefix:               strings.Join([]string{"hadoop", component, service}, "_"),
		urls:                 urls,
		commonLabels:         make(map[string]map[string][]string),
		firstGetCommonLabels: firstGetCommonLabels,
	}

	logrus.Infof("collector %s/%s - %v", collector.service, component, collector.urls)

	var allRules map[string][]MetricRulePattern
	if commonRules != nil {
		allRules = commonRules.MetricRules
	}
	if serviceRules != nil {
		allRules = utils.MergeMaps(allRules, serviceRules.MetricRules)
		collector.rules = allRules
		collector.isLowerName = serviceRules.LowercaseOutputName
		collector.isLowerLabel = serviceRules.LowercaseOutputLabel
	}
	collector.refreshMetricDescs()
	return collector
}

func (collector *MetricCollector) refreshMetricDescs() {
	metricDescs := make(map[string]map[string]*prometheus.Desc)
	collector.metricDescs = &metricDescs
}

func (collector *MetricCollector) getCommonLabels(beans *[]map[string]interface{}, url string) {
	collector.firstGetCommonLabels[url] = false
	collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "cluster")
	collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], collector.cluster)

	switch collector.service {
	case "namenode":
		beanPattern := "Hadoop:service=NameNode,name=JvmMetrics"
		bean := findBean(beans, beanPattern)
		if bean != nil {
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], (*bean)["tag.Hostname"].(string))
		}
	case "datanode":
		beanPattern := "Hadoop:service=DataNode,name=JvmMetrics"
		bean := findBean(beans, beanPattern)
		if bean != nil {
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], (*bean)["tag.Hostname"].(string))
		}
	case "journalnode":
		beanPattern := "Hadoop:service=JournalNode,name=JvmMetrics"
		bean := findBean(beans, beanPattern)
		if bean != nil {
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], (*bean)["tag.Hostname"].(string))
		}
	case "resourcemanager":
		beanPattern := "Hadoop:service=ResourceManager,name=JvmMetrics"
		bean := findBean(beans, beanPattern)
		if bean != nil {
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], (*bean)["tag.Hostname"].(string))
		}
	case "nodemanager":
		beanPattern := "Hadoop:service=NodeManager,name=JvmMetrics"
		bean := findBean(beans, beanPattern)
		if bean != nil {
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], (*bean)["tag.Hostname"].(string))
		}
	case "hiveserver2":
		beanPattern := `org.apache.logging.log4j2:type=AsyncContext@(\w{8})$`
		bean := findBean(beans, beanPattern)
		if bean != nil {
			subPattern := regexp2.MustCompile(".*hostName=(.+),.*", 0)
			matched, _ := subPattern.FindStringMatch((*bean)["ConfigProperties"].(string))
			collector.commonLabels[url]["names"] = append(collector.commonLabels[url]["names"], "host")
			collector.commonLabels[url]["values"] = append(collector.commonLabels[url]["values"], matched.Groups()[0].String())
		}
	default:
		return
	}
}

func (collector *MetricCollector) convertMetrics(beans *[]map[string]interface{}, url string, metrics *map[string]map[string][]prometheus.Metric) {
	for _, bean := range *beans {
		beanName, _ := bean["name"].(string)
		for beanPattern, metricRules := range collector.rules {
			beanPatternRegex := regexp2.MustCompile(beanPattern, 0)
			if ok, _ := beanPatternRegex.MatchString(beanName); !ok {
				continue
			}
			for beanMetricName, beanMetricValue := range bean {
				if utils.Contains(NON_METRIC_NAMES, beanMetricName) {
					continue
				}
				for _, metricDef := range metricRules {
					if metricDef.Type != "GAUSE" {
						logrus.Warnf("Only supporting GAUSE metric. Skip %s", metricDef.Type)
						continue
					}
					metricPatternRegex := regexp2.MustCompile(metricDef.Pattern, 0)
					if ok, _ := metricPatternRegex.MatchString(beanMetricName); ok {
						mergePattern := regexp2.MustCompile(fmt.Sprintf("%s<>%s", strings.TrimRight(beanPattern, "$"), strings.TrimLeft(metricDef.Pattern, "^")), 0)
						mergeName := fmt.Sprintf("%s<>%s", beanName, beanMetricName)
						subName, _ := mergePattern.Replace(mergeName, metricDef.Name, -1, -1)
						subLabelNames, subLabelValues := utils.ParseMap(metricDef.Labels)
						metricIdentifier := strings.Join(append([]string{subName}, subLabelNames...), "_")
						if _, ok := (*collector.metricDescs)[beanPattern][metricIdentifier]; !ok {
							prometheusMetricName := strings.Join([]string{collector.prefix, subName}, "_")
							if collector.isLowerName {
								prometheusMetricName = strings.ToLower(prometheusMetricName)
							}
							prometheusLabelNames := append(collector.commonLabels[url]["names"], subLabelNames...)
							if collector.isLowerLabel {
								prometheusLabelNames = utils.LambdaApply(prometheusLabelNames, strings.ToLower)
							}
							prometheusMetricHelp := ""
							if strings.Trim(metricDef.Help, " ") == "" {
								prometheusMetricHelp = prometheusMetricName
							} else {
								prometheusMetricHelp, _ = mergePattern.Replace(mergeName, metricDef.Help, -1, -1)
							}
							metricDesc := prometheus.NewDesc(prometheusMetricName, prometheusMetricHelp, prometheusLabelNames, nil)
							(*collector.metricDescs)[beanPattern][metricIdentifier] = metricDesc
							(*metrics)[beanPattern][metricIdentifier] = make([]prometheus.Metric, 0)
						}
						if metricDesc, ok := (*collector.metricDescs)[beanPattern][metricIdentifier]; ok {
							subLabelValues := utils.LambdaApply(
								subLabelValues,
								func(label string) string {
									result, _ := mergePattern.Replace(mergeName, label, -1, -1)
									return result
								},
							)
							prometheusLabelValues := append(collector.commonLabels[url]["values"], subLabelValues...)
							if collector.isLowerLabel {
								prometheusLabelValues = utils.LambdaApply(prometheusLabelValues, strings.ToLower)
							}
							prometheusMetricValue := resolveMetricValue(beanName, beanMetricName, beanMetricValue)
							prometheusMetric := prometheus.MustNewConstMetric(
								metricDesc,
								prometheus.GaugeValue,
								prometheusMetricValue,
								prometheusLabelValues...,
							)
							(*metrics)[beanPattern][metricIdentifier] = append((*metrics)[beanPattern][metricIdentifier], prometheusMetric)
						}
						break
					}
				}
			}
		}
	}
}

func (collector MetricCollector) Collect(ch chan<- prometheus.Metric) {
	collector.refreshMetricDescs()
	metrics := make(map[string]map[string][]prometheus.Metric)

	for beanPattern := range collector.rules {
		(*collector.metricDescs)[beanPattern] = make(map[string]*prometheus.Desc)
		metrics[beanPattern] = make(map[string][]prometheus.Metric)
	}

	for _, url := range *collector.urls {
		beans := getRawBeans(url)
		if beans == nil {
			logrus.Warnf("Can't scrape metrics from %s", url)
			continue
		}
		if val, ok := collector.firstGetCommonLabels[url]; ok && val {
			collector.commonLabels[url] = map[string][]string{
				"names":  make([]string, 0),
				"values": make([]string, 0),
			}
			collector.getCommonLabels(beans, url)
		}
		collector.convertMetrics(beans, url, &metrics)
	}
	for _, beanMetrics := range metrics {
		for _, metrics := range beanMetrics {
			for _, metric := range metrics {
				ch <- metric
			}
		}
	}
}

func (collector MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, beanMetricDescs := range *collector.metricDescs {
		for _, metricDesc := range beanMetricDescs {
			ch <- metricDesc
		}
	}
}

func resolveMetricValue(beanName string, name string, value interface{}) float64 {
	if strings.Contains(beanName, "FSNamesystemState") && strings.Contains(name, "FSState") {
		return mapFSState(value.(string))
	} else if strings.Contains(beanName, "FSNamesystem") && strings.Contains(name, "HAState") {
		return mapHAState(value.(string))
	} else if strings.Contains(beanName, "RMInfo") && name == "State" {
		return mapRMState(value.(string))
	}
	return value.(float64)
}

func mapFSState(v string) float64 {
	switch v {
	case "Operational":
		return 0
	case "Safemode":
		return 1
	default:
		return 9999
	}
}

func mapHAState(v string) float64 {
	switch v {
	case "initializing":
		return 0
	case "active":
		return 1
	case "standby":
		return 2
	case "stopping":
		return 3
	default:
		return 9999
	}
}

func mapRMState(v string) float64 {
	return mapHAState(v)
}

func getRawBeans(url string) *[]map[string]interface{} {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Warnf("Error when getMetric: %v", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("Get %s failed, response code is: %d", url, resp.StatusCode)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Warnf("Read response body failded: %v", err)
		return nil
	}
	var data map[string][]map[string]interface{}
	_ = json.Unmarshal(body, &data)
	if val, ok := data["beans"]; ok {
		return &val
	}
	logrus.Warnf("Got zone metrics in from %s", url)
	return nil
}

func findBean(beans *[]map[string]interface{}, beanPattern string) *map[string]interface{} {
	pattern := regexp2.MustCompile(beanPattern, 0)
	for _, bean := range *beans {
		if ok, _ := pattern.MatchString(bean["name"].(string)); ok {
			return &bean
		}
	}
	return nil
}

type TestCollector struct {
	metricDesc []*prometheus.Desc
	metrics    []prometheus.Metric
}

func InitTestCollector() TestCollector {
	delta := rand.Intn(1000)
	desc := make([]*prometheus.Desc, 10)
	for i := 0; i < 10; i++ {
		desc[i] = prometheus.NewDesc(
			fmt.Sprintf("Name%d", i),
			fmt.Sprintf("Help %d", i),
			[]string{"step"},
			nil,
		)
	}
	metrics := make([]prometheus.Metric, 10)
	for i := 0; i < 10; i++ {
		metrics[i] = prometheus.MustNewConstMetric(desc[i], prometheus.GaugeValue, float64(i*delta), fmt.Sprint(i*delta))
	}
	return TestCollector{
		metricDesc: desc,
		metrics:    metrics,
	}
}

func (collector TestCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range collector.metricDesc {
		ch <- desc
	}
}

func (collector TestCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range collector.metrics {
		ch <- metric
	}
}
