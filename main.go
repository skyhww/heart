package main

import (
	"github.com/astaxie/beego"
	"heart/controller"
	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
	"heart/cfg"
)

func main() {
	beego.Router("/token", &controller.Token{})
	beego.Run()
}

var redisPool *redis.Pool

func init() {
	cfg := loadConfig()
	if cfg == nil {
		panic(errors.New("加载配置文件失败！"))
	}
	//loadRedis
	loadRedis(cfg.RedisConfig)

}

func loadConfig() *cfg.Config {
	c := &cfg.Config{RedisConfig: &cfg.RedisConfig{}}
	cfg.LoadYaml("etc/redis.yml", c.RedisConfig)
	return c
}

func loadRedis(redisConfig *cfg.RedisConfig) {
	redisPool = &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		IdleTimeout: time.Duration(redisConfig.IdleTimeOutSeconds) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisConfig.Url)
			if err != nil {
				panic(err)
			}
			if _, authErr := c.Do("AUTH", redisConfig.Auth); authErr != nil {
				panic(authErr)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
