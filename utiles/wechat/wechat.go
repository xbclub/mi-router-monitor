package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"net/http"
)

const tokenKey = "mirouterMonitor:wechat:token:%s"

type Wechatconfig struct {
	corpid     string
	corpsecret string
	appid      int
	siteproxy  string
	redis      *redis.Redis
}
type Sendmails interface {
	Sendmail(msg string) bool
}

func NewWechat(corpid, corpsecret, siteproxy string, appid int, rediscon *redis.Redis) Sendmails {
	return &Wechatconfig{
		corpid:     corpid,
		corpsecret: corpsecret,
		appid:      appid,
		redis:      rediscon,
		siteproxy:  siteproxy,
	}
}

//	func Do(method string, url string, payload io.Reader) (*http.Response, error) {
//		req, err := http.NewRequest(method, url, payload)
//		if err != nil {
//			return nil, err
//		}
//		return http.DefaultClient.Do(req)
//	}
func (c *Wechatconfig) getAcessKey() (res string, err error) {
	get, _ := c.redis.Get(fmt.Sprintf(tokenKey, c.corpid))
	if len(get) > 0 {
		return get, nil
	}
	logx.Info("token 缓存不存在 重新申请token")
	url := "%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	do, err := httpc.Do(context.Background(), http.MethodGet, fmt.Sprintf(url, c.siteproxy, c.corpid, c.corpsecret), nil)
	if err != nil {
		return "", err
	}
	//ress, _ := Do("GET", fmt.Sprintf(url, c.corpid, c.corpsecret), nil)
	if do.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("wexin get token http request error status code %d", do.StatusCode))
	}
	defer do.Body.Close()
	accesstoken, err := tojson(do)
	if err != nil || !errorhandle(accesstoken) {
		logx.Errorf(err.Error())
		return res, errors.New("get access token response error")
	}
	c.redis.SetnxEx(fmt.Sprintf(tokenKey, c.corpid), accesstoken["access_token"].(string), 7100)
	return accesstoken["access_token"].(string), nil
}
func (c *Wechatconfig) Sendmail(msg string) bool {
	url := "%s/cgi-bin/message/send?access_token=%s"
	accesstoken, err := c.getAcessKey()
	if err != nil {
		logx.Errorf(err.Error())
		return false
	}
	body := SendText{
		Touser:  "@all",
		Msgtype: "text",
		Agentid: c.appid,
		Text: Text{
			Content: msg,
		},
		Safe:                     0,
		Enable_duplicate_check:   0,
		Duplicate_check_interval: 0,
	}
	res, err := httpc.Do(context.Background(), "POST", fmt.Sprintf(url, c.siteproxy, accesstoken), body)
	if err != nil {
		logx.Error(err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		logx.Errorf("wexin send msg http request error status code %d", res.StatusCode)
		return false
	}
	jsons, _ := tojson(res)
	status := errorhandle(jsons)
	if !status {
		logx.WithContext(context.Background()).Errorf(jsons["errmsg"].(string))
		_, err := c.redis.Del(fmt.Sprintf(tokenKey, c.corpid))
		if err != nil {
			logx.Error("请求删除微信失败token错误")
		}
		return false
	}
	return status
}
func tojson(res *http.Response) (jsons map[string]interface{}, err error) {
	if res == nil {
		return jsons, errors.New("请求未获取到数据")
	}
	myres, _ := io.ReadAll(res.Body)
	jsons = make(map[string]interface{})
	//fmt.Printf("res:%s",string(myres))
	err = json.Unmarshal(myres, &jsons)
	return
}
func errorhandle(jsons map[string]interface{}) bool {
	if jsons["errcode"].(float64) != 0 {
		logx.Errorf("get token error: %s", jsons["errmsg"].(string))
		return false
	} else {
		return true
	}
}
