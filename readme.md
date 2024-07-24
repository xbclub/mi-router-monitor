### 本项目概述

本项目基于 [mirouter-ui](https://github.com/Mirouterui/mirouter-ui) 接口构建的上下行网络速率告警

### 功能列表

- [x] 上行速率告警
- [x] 下行速率告警
- [x] 告警频率控制
- [x] 企业微信应用推送

### 项目依赖

#### 1. 启动依赖

- Mirouter-ui: 1.3.6+
- Redis: 7.2.4+

#### 2. 编译依赖

- GO: 1.22.5+

### Docker-Compose

#### 1. 环境变量配置

| 变量名                          | 默认值                         | 描述                    |
|------------------------------|-----------------------------|-----------------------|
| MIROUTER_UI_URL              | 无                           | mirouter-ui 地址        |
| REDIS_HOST                   | 无                           | redis 地址              |
| REDIS_TYPE                   | node                        | redis 类型              |
| REDIS_PASS                   | 空                           | redis 密码 `没有不要添加这个变量` |
| REDIS_TLS                    | false                       | redis 是否开启tls         |
| WECHAT_CORPID                | 无                           | 企业微信 企业秘钥             |
| WECHAT_CORPSECRET            | 无                           | 企业微信 应用秘钥             |
| WECHAT_APPID                 | 无                           | 企业微信 应用ID             |
| WECHAT_PROXYSITE             | https://qyapi.weixin.qq.com | 企业微信 api接口域名          |
| MONITOR_UPLOAD_SPEED_LIMIT   | 3145728                     | 上行告警速率阈值 单位B/s        |
| MONITOR_DOWNLOAD_SPEED_LIMIT | 31457280                    | 下行告警速率阈值 单位B/s        |
| MONITOR_ALERT_QUOTA          | 600                         | 重复告警间隔 单位 s           |
#### 2. 启动命令
```bash
# 创建文件夹
mkdir mirouter-monitor && cd mirouter-monitor
# 下载 docker-compose 文件
wget https://raw.githubusercontent.com/xbclub/mi-router-monitor/main/docker-compose.yml
# 启动项目
docker compose up -d
# 配置 mirouter-ui
vim ./data/mirouter-ui/data/config.json
# 重启服务
docker compose restart
```

### 鸣谢

- [mirouter-ui](https://github.com/Mirouterui/mirouter-ui)