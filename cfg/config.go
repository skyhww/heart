package cfg

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Config struct {
	RedisConfig    *RedisConfig
	AliYunConfig   *AliYunConfig
	DatabaseConfig *DatabaseConfig
	ElasticConfig  *ElasticConfig
}

//file: location
//r: decode config
func LoadYaml(file string, r interface{}) {
	fmt.Printf("读取配置文件%s开始\n", file)
	defer fmt.Printf("配置文件%s加载完成\n", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		panic(err)
	}
}
