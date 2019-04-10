// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package main


import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
	"heart/cfg"
	_ "github.com/go-sql-driver/mysql"
	"heart/sms"
	"github.com/jmoiron/sqlx"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"heart/service"
	"heart/entity"
	"heart/controller"
	"heart/controller/common"
)

func main() {
	//读取配置文件
	cfg := loadConfig()
	//加载配置  阿里云、redis、数据库
	aliYun := loadAliyun(cfg.AliYunConfig)
	redisPool := loadRedis(cfg.RedisConfig)
	db := loadDatabase(cfg.DatabaseConfig)

	//helper
	tokenHelper:=&service.TokenHelper{Rds:redisPool}
	simpleTokenService:=&service.SimpleTokenService{Pool:redisPool,Ex: time.Second*60*60*3}
	//persist
	userPersist:=entity.NewUserPersist(db)
	storePersist:=entity.NewStorePersist(db)
	userVideoPersist:=entity.NewUserVideoPersist(db)
	//security
	security := &service.SimpleSecurity{Pool: redisPool, SmsClient: aliYun, UserPersist:userPersist,TokenService:simpleTokenService}
	//service
	storeService:= &service.LocalStoreService{Path:"store",Type:"LOCAL",StorePersist:storePersist}
	userInfo:=&service.UserInfo{UserPersist:userPersist,StoreService:storeService}
	videoService:=&service.SimpleVideoService{StoreService:storeService,UserPersist:userPersist,UserVideoPersist:userVideoPersist}
	//SmsController
	smsController := &controller.SmsController{}
	smsController.Security = security
	//Token
	token := &controller.Token{Service: security}
	tokenHolder:=&common.TokenHolder{Name:"token",Helper:tokenHelper}
	//user
	userController := &controller.User{Service: security,TokenHolder:tokenHolder}
	userName:=&controller.Name{TokenHolder:tokenHolder,UserInfo:userInfo}
	signature:=&controller.Signature{TokenHolder:tokenHolder,UserInfo:userInfo}
	icon:=&controller.Icon{TokenHolder:tokenHolder,UserInfo:userInfo,Limit:3}
	videoController:= &controller.VideoController{VideoService:videoService,TokenHolder:tokenHolder,Limit:50}
	//iMessage
	iMessage:=&controller.IMessageController{}
	//beego运行
	ns := beego.NewNamespace("/heart/v1.0")
	ns.Router("/token", token)
	ns.Router("/sms/:mobile", smsController)
	ns.Router("/user", userController)
	ns.Router("/user/info/name",userName)
	ns.Router("/user/info/signature",signature)
	ns.Router("/user/info/icon",icon)
	ns.Router("/video",videoController)
	ns.Handler("/message",iMessage)
	beego.AddNamespace(ns)
	beego.Run()
}

func loadConfig() *cfg.Config {
	c := &cfg.Config{RedisConfig: &cfg.RedisConfig{}, AliYunConfig: &cfg.AliYunConfig{}, DatabaseConfig: &cfg.DatabaseConfig{}}
	cfg.LoadYaml("conf/redis.yml", c.RedisConfig)
	cfg.LoadYaml("conf/aliyun.yml", c.AliYunConfig)
	cfg.LoadYaml("conf/database.yml", c.DatabaseConfig)
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
