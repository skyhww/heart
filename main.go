package main

import (
	"github.com/astaxie/beego"
	"heart/controller"
	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
	"heart/cfg"
	_ "github.com/go-sql-driver/mysql"
	"heart/sms"
	"github.com/jmoiron/sqlx"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
)

func main() {
	//

	ns := beego.NewNamespace("/v1.0")
	ns.Router("/token", &controller.Token{})
	ns.Router("/sms/:mobile", &controller.SmsController{})
	beego.Run()
}

func init() {
	cfg := loadConfig()
	if cfg == nil {
		panic(errors.New("加载配置文件失败！"))
	}
}

func loadConfig() *cfg.Config {
	c := &cfg.Config{RedisConfig: &cfg.RedisConfig{}, AliYunConfig: &cfg.AliYunConfig{}, DatabaseConfig: &cfg.DatabaseConfig{}}
	cfg.LoadYaml("etc/redis.yml", c.RedisConfig)
	cfg.LoadYaml("etc/aliyun.yml", c.AliYunConfig)
	cfg.LoadYaml("etc/database.yml", c.DatabaseConfig)
	return c
}

func loadAliyun(aliYunConfig *cfg.AliYunConfig) *sms.AliYun {
	client, err := sdk.NewClientWithAccessKey(aliYunConfig.RegionId, aliYunConfig.AccessKeyId, aliYunConfig.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	return &sms.AliYun{AliYunConfig: aliYunConfig, Client: client}
}

func loadDatabase(databaseConfig *cfg.DatabaseConfig) *sqlx.DB {
	db, err := sqlx.Open(databaseConfig.Driver, databaseConfig.Dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(databaseConfig.MaxOpenConns)
	db.SetMaxIdleConns(databaseConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(databaseConfig.ConnMaxLifetimeSeconds) * time.Second)
	return db
}
func loadRedis(redisConfig *cfg.RedisConfig) *redis.Pool {
	return &redis.Pool{
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
