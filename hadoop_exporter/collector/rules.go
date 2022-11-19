package collector

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type MetricRulePattern struct {
	Pattern string            `yaml:"pattern"`
	Type    string            `yaml:"type"`
	Name    string            `yaml:"name"`
	Labels  map[string]string `yaml:"labels"`
	Help    string            `yaml:"help"`
}

type ServiceRules struct {
	LowercaseOutputName  bool                           `yaml:"lowercaseOutputName"`
	LowercaseOutputLabel bool                           `yaml:"lowercaseOutputLabel"`
	MetricRules          map[string][]MetricRulePattern `yaml:"rules"`
}

func ReadServiceRules(path string) (*ServiceRules, error) {
	buffer, err := ioutil.ReadFile(path)
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
