package utils

import (
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v2"
)

func ReadYamlFile(path string) (map[string]interface{}, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		logrus.Warnf("Error when attempt to read %s: %v", path, err)
		return nil, err
	}
	var data map[string]interface{}
	err = yaml.Unmarshal(buffer, &data)
	if err != nil {
		logrus.Warnf("Error when attemp to parse %s: %v", string(buffer[:]), err)
		return nil, err
	}
	return data, nil
}

func Handlelogger(logLevel string) {
	parsedLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Errorf("Invalid log level: %s. Using the default log level: info", logLevel)
	} else {
		logrus.SetLevel(parsedLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:             false,
		FullTimestamp:             true,
		EnvironmentOverrideColors: true,
	})
}

func Coalesce(values ...*interface{}) *interface{} {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}

func CoalesceString(values ...string) string {
	for _, v := range values {
		if strings.Trim(v, " ") != "" {
			return v
		}
	}
	return ""
}

func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func KeysOf[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ValueOf[K constraints.Ordered, V any](m map[K]V) []V {
	values := make([]V, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func ParseMap[K constraints.Ordered, V any](m map[K]V) ([]K, []V) {
	keys := KeysOf(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	values := make([]V, len(m))
	for i, k := range keys {
		values[i] = m[k]
	}
	return keys, values
}

func LambdaApply[T any, V any](slice []T, apply func(T) V) []V {
	result := make([]V, len(slice))
	for i, val := range slice {
		result[i] = apply(val)
	}
	return result
}

func MergeMaps[K comparable, V any](a map[K]V, b map[K]V) map[K]V {
	r := a
	for k, v := range b {
		r[k] = v
	}
	return r
}
