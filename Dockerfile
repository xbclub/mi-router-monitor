from alpine
run sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk --no-cache add tzdata
add app /app/
add etc/docker-config.yaml /app/etc/config.yaml
workdir /app
ENV REDIS_TYPE="node" \
    REDIS_PASS="" \
    REDIS_TLS=false \
    WECHAT_PROXYSITE="https://qyapi.weixin.qq.com" \
    MONITOR_UPLOAD_SPEED_LIMIT=3145728 \
    MONITOR_DOWNLOAD_SPEED_LIMIT=31457280 \
    MONITOR_ALERT_QUOTA=600
ENTRYPOINT ["./app"]