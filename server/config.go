package server

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type HadoopClusterJmx struct {
	Services map[string]*[]string `yaml:"services"`
	Cluster  string               `yaml:"cluster"`
}

type ExporterConfig struct {
	Server *struct {
		Address    string `yaml:"address"`
		MetricPath string `yaml:"metricPath"`
		LogLevel   string `yaml:"logLevel"`
		Port       int    `yaml:"port"`
	} `yaml:"server"`
	Jmx *[]HadoopClusterJmx `yaml:"jmx"`
}

func ReadExporterConfig(path string) (*ExporterConfig, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		logrus.Warnf("Error when attempt to read %s: %v", path, err)
		return nil, err
	}
	var data ExporterConfig
	err = yaml.Unmarshal(buffer, &data)
	if err != nil {
		logrus.Warnf("Error when attemp to parse %s: %v", string(buffer[:]), err)
		return nil, err
	}
	return &data, nil
}
