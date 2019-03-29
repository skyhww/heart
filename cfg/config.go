package cfg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RedisConfig    *RedisConfig
	AliYunConfig   *AliYunConfig
	DatabaseConfig *DatabaseConfig
}

//file: location
//r: decode config
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
