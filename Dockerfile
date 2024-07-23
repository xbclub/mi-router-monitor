from alpine
run apk --no-cache add tzdata
add app /app/
workdir /app
ENTRYPOINT ["./app"]