package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Discovery struct {
	Conf Config `yaml:"discovery"`
}

type Config struct {
	ConfVersion    string            `yaml:"version"`
	Address        string            `yaml:"address"`
	Db             []string          `yaml:"db"`
	Heartbeat      string            `yaml:"heartbeat"`
	HeartbeatTopic string            `yaml:"heartbeatTopic"`
	InstrumentConf map[string]string `yaml:"instrument"`
}

func ConfigFile(n ...string) (*Config, error) {
	path := "config.yml"
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	check(err)

	var conf Discovery
	err = yaml.Unmarshal(yamlFile, &conf)
	check(err)

	return &conf.Conf, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
