package main

import (
	"github.com/astaxie/beego"
	"heart/controller"
	"github.com/garyburd/redigo/redis"
	"time"
	"heart/cfg"
	_ "github.com/go-sql-driver/mysql"
	"heart/sms"
	"github.com/jmoiron/sqlx"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"heart/service"
	"heart/entity"
)

func main() {
	//读取配置文件
	cfg := loadConfig()
	//加载配置  阿里云、redis、数据库
	aliYun := loadAliyun(cfg.AliYunConfig)
	redisPool := loadRedis(cfg.RedisConfig)
	db := loadDatabase(cfg.DatabaseConfig)
	//service
	service := &service.SimpleSecurity{Pool: redisPool, SmsClient: aliYun, UserPersist: entity.NewUserPersist(db)}
	//SmsController
	smsController := &controller.SmsController{}
	smsController.Security = service
	//Token
	token := &controller.Token{Service: service}
	//user
	userController := &controller.User{Service: service}
	//beego运行
	ns := beego.NewNamespace("/v1.0")
	ns.Router("/token", token)
	ns.Router("/sms/:mobile", smsController)
	ns.Router("/user", userController)
	beego.Run()
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
