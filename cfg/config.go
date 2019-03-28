package cfg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type RedisConfig struct {
	Url                string `yaml:"url"`
	Auth               string `yaml:"auth"`
	MaxIdle            int    `yaml:"max_idle"`
	IdleTimeOutSeconds int    `yaml:"idle_time_out_seconds"`
}

type Config struct {
	RedisConfig *RedisConfig
}

func LoadYaml(file string, r interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		panic(err)
	}
}
