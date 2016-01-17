package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OpenHABHost   string
	WifiDevice    string `yaml:"wifidevice,omitempty"`
	MonitorDevice string `yaml:"monitordevice,omitempty"`
	MonitoredMACs map[string]string
}

var GlobalConfig Config

func ParseConfig(configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Can't read configuration file %s:%v", configPath, err)
	}
	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Can't unmarshal config: %v", err)
	}
	GlobalConfig = config
}
