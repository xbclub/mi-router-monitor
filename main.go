package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"mirouterMoinitor/config"
	"mirouterMoinitor/svc"
	"mirouterMoinitor/utiles/alert"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	// 允许配置文件调用环境变量
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	// 关闭指标日志输出
	logx.DisableStat()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	running := true
	for running {
		select {
		case <-sig:
			running = false
			logx.Info("接收到程序退出信号，程序退出")
		case <-t.C:
			miroutertmp := alert.NewMirouterConnect(ctx)
			status, err := miroutertmp.GetStatus()
			if err != nil {
				logx.Error(err)
				t.Reset(interval)
				continue
			}
			miroutertmp.ComputeUploadSpeed(status)
			t.Reset(interval)
		}
	}

}
