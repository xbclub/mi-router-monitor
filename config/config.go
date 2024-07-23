package config

import "github.com/zeromicro/go-zero/core/stores/redis"

type Config struct {
	Redis        redis.RedisConf
	MiRouterURL  string
	Wechatconfig struct {
		Corpid     string
		Corpsecret string
		Appid      int
		ProxySite  string `json:",default=https://qyapi.weixin.qq.com"`
	}
	MonitorConf struct {
		UploadSpeedLimit   int64 `json:",default=3072"`
		DownloadSpeedLimit int64 `json:",default=3072"`
		// 告警间隔 单位秒
		AlertQuota int `json:",default=600"`
	}
}
