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
	"github.com/astaxie/beego/plugins/cors"
	"heart/extend"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}
func main() {
	//读取配置文件
	cfg := loadConfig()
	//加载配置  阿里云、redis、数据库
	aliYun := loadAliyun(cfg.AliYunConfig)
	redisPool := loadRedis(cfg.RedisConfig)
	db := loadDatabase(cfg.DatabaseConfig)
	s := extend.NewSegmentsFilter()
	//helper
	tokenHelper := &service.TokenHelper{Rds: redisPool}
	simpleTokenService := &service.SimpleTokenService{Pool: redisPool, Ex: time.Second * 60 * 60 * 24}
	//persist
	userPersist := entity.NewUserPersist(db)
	storePersist := entity.NewStorePersist(db)
	userVideoPersist := entity.NewUserVideoPersist(db)
	messagePersist := entity.NewMessagePersist(db)
	postCommentPersist := entity.NewPostCommentPersist(db)
	postAttachPersist := entity.NewPostAttachPersist(db)
	userPostPersist := entity.NewUserPostPersist(db)
	postsPersist := entity.NewPostsPersist(db)
	userFollowInfoPersist := entity.NewUserFollowInfoPersist(db)
	userCollectionInfoPersist := entity.NewUserCollectionInfoPersist(db)
	//security
	security := &service.SimpleSecurity{Pool: redisPool, SmsClient: aliYun, UserPersist: userPersist, TokenService: simpleTokenService}
	//service
	elasticSearchService, err := service.NewElasticSearchService(cfg.ElasticConfig.Host)
	if err != nil {
		panic(err)
	}
	storeService := &service.LocalStoreService{Path: "store", Type: "LOCAL", StorePersist: storePersist}
	userInfo := &service.UserInfo{UserPersist: userPersist, StoreService: storeService}
	videoService := &service.SimpleVideoService{ElasticSearchService: elasticSearchService, StoreService: storeService, UserPersist: userPersist, UserVideoPersist: userVideoPersist}
	messageService := &service.SimpleMessageService{MessagePersist: messagePersist, UserPersist: userPersist, StoreService: storeService}
	userPostService := &service.SimpleUserPostService{SegmentsFilter: s, PostCommentPersist: postCommentPersist, PostAttachPersist: postAttachPersist, UserPersist: userPersist, UserPostPersist: userPostPersist}
	postAttachService := &service.SimplePostAttachService{PostAttachPersist: postAttachPersist, StoreService: storeService}
	postService := &service.SimplePostService{PostsPersist: postsPersist, PostAttachPersist: postAttachPersist, ElasticSearchService: elasticSearchService}
	userFollowService := &service.SimpleUserFollowService{UserPersist: userPersist, UserFollowInfoPersist: userFollowInfoPersist}
	collectorService := &service.SimpleCollectorService{UserCollectionInfoPersist: userCollectionInfoPersist, UserPersist: userPersist}
	//SmsController
	smsController := &controller.SmsController{}
	smsController.Security = security
	//Token
	token := &controller.Token{Service: security}
	tokenHolder := &common.TokenHolder{Name: "token", Helper: tokenHelper}
	//user
	userController := &controller.User{Service: security, TokenHolder: tokenHolder}
	userName := &controller.Name{TokenHolder: tokenHolder, UserInfo: userInfo}
	signature := &controller.Signature{TokenHolder: tokenHolder, UserInfo: userInfo}
	icon := &controller.Icon{TokenHolder: tokenHolder, UserInfo: userInfo, Limit: 3}
	videoController := &controller.VideoController{VideoService: videoService, TokenHolder: tokenHolder, Limit: 50}
	userPostsController := &controller.UserPostsController{UserPostService: userPostService, TokenHolder: tokenHolder, StoreService: storeService, Limit: 10, MaxAttach: 10}
	postAttachController := &controller.PostAttachController{TokenHolder: tokenHolder, PostAttachService: postAttachService}
	postsController := &controller.PostsController{PostService: postService, TokenHolder: tokenHolder}
	relationController := &controller.RelationController{TokenHolder: tokenHolder, UserFollowService: userFollowService}
	userCollectorController := &controller.UserCollectorController{TokenHolder: tokenHolder, CollectorService: collectorService}
	commentController := &controller.CommentController{PostService: userPostService, TokenHolder: tokenHolder}
	//iMessage
	iMessage := &controller.IMessageController{MessageService: messageService, TokenHolder: tokenHolder, Limit: 10}
	iMessageAttachController := &controller.IMessageAttachController{TokenHolder: tokenHolder, MessageService: messageService, Limit: 5}
	//beego运行
	ns := beego.NewNamespace("/heart/v1.0")
	ns.Router("/token", token)
	ns.Router("/sms/:mobile", smsController)
	ns.Router("/user", userController)
	ns.Router("/user/info/name", userName)
	ns.Router("/user/info/signature", signature)
	ns.Router("/message", iMessage)
	ns.Router("/message/:id/attach", iMessageAttachController)

	ns.Router("/user/posts", userPostsController)

	ns.Router("/user/posts/:id", userPostsController, "delete:Delete")
	ns.Router("/user/:id/icon", icon)

	ns.Router("/video", videoController, "get:Search")
	ns.Router("/video/:id", videoController)
	ns.Router("/video", videoController,"put:Put")

	ns.Router("/posts", postsController)
	ns.Router("/posts/attach/:id", postAttachController, "get:Get")
	ns.Router("/posts/:posts_id/attach", postAttachController, "get:GetPage")

	ns.Router("/comment/:id/comment", commentController, "put:Replay")
	ns.Router("/posts/:post_id/comment", commentController)
	ns.Router("/comment/:id", commentController, "delete:Delete")
	ns.Router("/relation/:user_id", relationController)
	ns.Router("/posts_collector/:posts_id", userCollectorController)

	beego.AddNamespace(ns)
	beego.Run()
}

func loadConfig() *cfg.Config {
	c := &cfg.Config{RedisConfig: &cfg.RedisConfig{}, AliYunConfig: &cfg.AliYunConfig{}, DatabaseConfig: &cfg.DatabaseConfig{}, ElasticConfig: &cfg.ElasticConfig{}}
	cfg.LoadYaml("conf/redis.yml", c.RedisConfig)
	cfg.LoadYaml("conf/aliyun.yml", c.AliYunConfig)
	cfg.LoadYaml("conf/database.yml", c.DatabaseConfig)
	cfg.LoadYaml("conf/elastic.yml", c.ElasticConfig)
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
