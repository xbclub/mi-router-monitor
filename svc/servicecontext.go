package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"mirouterMoinitor/config"
	"mirouterMoinitor/utiles/wechat"
)

type ServiceContext struct {
	Config config.Config
	RedisC *redis.Redis
	Wechat wechat.Sendmails
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisC := redis.MustNewRedis(c.Redis)
	return &ServiceContext{
		Config: c,
		RedisC: redisC,
		Wechat: wechat.NewWechat(c.Wechatconfig.Corpid, c.Wechatconfig.Corpsecret, c.Wechatconfig.ProxySite, c.Wechatconfig.Appid, redisC),
	}
}
