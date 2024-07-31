package alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"mirouterMoinitor/svc"
	"net/http"
	"strconv"
)

type MirouterConnect struct {
	svc             *svc.ServiceContext
	uploadlimiter   *limit.PeriodLimit
	downloadlimiter *limit.PeriodLimit
}
type MirouterConnecter interface {
	GetStatus() (*MiRouterStatus, error)
	ComputeUploadSpeed(status *MiRouterStatus)
}

var limitkey = "miRouter:limit"
var UpSpeedOverLimitTemplate = `路由器告警
当前%s速度大于 %s
当前速度：%s
`
var deviceList = `
%v. %s 
%s速度: %v
`

func NewMirouterConnect(svc *svc.ServiceContext) MirouterConnecter {
	return &MirouterConnect{
		svc:             svc,
		uploadlimiter:   limit.NewPeriodLimit(svc.Config.MonitorConf.AlertQuota, 1, svc.RedisC, limitkey),
		downloadlimiter: limit.NewPeriodLimit(svc.Config.MonitorConf.AlertQuota, 1, svc.RedisC, limitkey),
	}
}
func (m *MirouterConnect) GetStatus() (*MiRouterStatus, error) {
	do, err := httpc.Do(context.Background(), http.MethodGet, m.svc.Config.MiRouterURL+"/0/api/misystem/status", nil)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http request error status code %d", do.StatusCode))
	}
	status := &MiRouterStatus{}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(all, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// 计算当前上传速度
func (m *MirouterConnect) ComputeUploadSpeed(status *MiRouterStatus) {
	if status.Wan.Upspeed >= m.svc.Config.MonitorConf.UploadSpeedLimit {
		result, err := m.uploadlimiter.Take("upload")
		if err != nil {
			logx.Errorf("Error: %v", err)
		}
		var tmpDeviceList = ""
		for i, x := range sortDevUp(status.Dev) {
			tmpDeviceList += fmt.Sprintf(deviceList, i+1, x.Devname, "上传", byteConvert(convertint64(x.Upspeed)))
		}
		if result == limit.Allowed || result == limit.HitQuota {
			wxstatus := m.svc.Wechat.Sendmail(fmt.Sprintf(UpSpeedOverLimitTemplate, "上传", byteConvert(m.svc.Config.MonitorConf.UploadSpeedLimit)+tmpDeviceList, byteConvert(status.Wan.Upspeed)))
			if wxstatus != true {
				logx.Error("微信推送失败"+UpSpeedOverLimitTemplate, "上传", byteConvert(m.svc.Config.MonitorConf.UploadSpeedLimit), byteConvert(status.Wan.Upspeed))
			}
		} else {
			logx.Info("上传告警指定时间段内已推送过，跳过")
		}
	}
	if status.Wan.Downspeed >= m.svc.Config.MonitorConf.DownloadSpeedLimit {
		result, err := m.downloadlimiter.Take("download")
		if err != nil {
			logx.Errorf("Error: %v", err)
		}
		var tmpDeviceList = ""
		for i, x := range sortDevDown(status.Dev) {
			tmpDeviceList += fmt.Sprintf(deviceList, i+1, x.Devname, "下载", byteConvert(convertint64(x.Upspeed)))
		}
		if result == limit.Allowed || result == limit.HitQuota {
			wxstatus := m.svc.Wechat.Sendmail(fmt.Sprintf(UpSpeedOverLimitTemplate, "下载", byteConvert(m.svc.Config.MonitorConf.DownloadSpeedLimit)+tmpDeviceList, byteConvert(status.Wan.Downspeed)))
			if wxstatus != true {
				logx.Error("微信推送失败"+UpSpeedOverLimitTemplate, "下载", byteConvert(m.svc.Config.MonitorConf.DownloadSpeedLimit), byteConvert(status.Wan.Downspeed))
			}
		} else {
			logx.Info("下载告警指定时间段内已推送过，跳过")
		}
	}

}

// 对 status.dev Upspeed字段进行排序
func sortDevUp(dev []Dev) []Dev {
	for i := 0; i < len(dev); i++ {
		for j := i + 1; j < len(dev); j++ {
			if convertint64(dev[i].Upspeed) < convertint64(dev[j].Upspeed) {
				dev[i], dev[j] = dev[j], dev[i]
			}
		}
	}
	return dev
}

// 对 status.dev Downspeed字段进行排序
func sortDevDown(dev []Dev) []Dev {
	for i := 0; i < len(dev); i++ {
		for j := i + 1; j < len(dev); j++ {
			if convertint64(dev[i].Downspeed) < convertint64(dev[j].Downspeed) {
				dev[i], dev[j] = dev[j], dev[i]
			}
		}
	}
	return dev
}
func convertint64(s interface{}) int64 {
	switch s.(type) {
	case string:
		parseInt, err := strconv.ParseInt(s.(string), 10, 64)
		if err != nil {
			logx.Error(err)
			return 0
		}
		return parseInt
	case int:
		return int64(s.(int))
	case int64:
		return s.(int64)
	default:
		return 0
	}
}

// 字节转换 kb mb gb
func byteConvert(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B/s", b)
	} else if b < 1024*1024 {
		return fmt.Sprintf("%.2f KB/s", float64(b)/1024)
	} else if b < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB/s", float64(b)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB/s", float64(b)/(1024*1024*1024))
	}
}
