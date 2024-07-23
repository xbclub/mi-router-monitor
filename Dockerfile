from alpine
run sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk --no-cache add tzdata
add app /app/
workdir /app
ENTRYPOINT ["./app"]