package collector

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type MetricRulePattern struct {
	Labels  map[string]string `yaml:"labels"`
	Pattern string            `yaml:"pattern"`
	Type    string            `yaml:"type"`
	Name    string            `yaml:"name"`
	Help    string            `yaml:"help"`
}

type ServiceRules struct {
	MetricRules          map[string][]MetricRulePattern `yaml:"rules"`
	LowercaseOutputName  bool                           `yaml:"lowercaseOutputName"`
	LowercaseOutputLabel bool                           `yaml:"lowercaseOutputLabel"`
}

func ReadServiceRules(path string) (*ServiceRules, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		logrus.Warnf("Error when attempt to read %s: %v", path, err)
		return nil, err
	}
	var data ServiceRules
	err = yaml.Unmarshal(buffer, &data)
	if err != nil {
		logrus.Warnf("Error when attemp to parse %s: %v", string(buffer[:]), err)
		return nil, err
	}
	return &data, nil
}
