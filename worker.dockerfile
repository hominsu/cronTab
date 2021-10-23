#build stage
FROM golang:alpine AS builder
# change apk source to ali
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk upgrade
# change go proxy to qiniu cdn
RUN go env -w GOPROXY=https://goproxy.cn,direct
# build app
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./worker/main/main.go

#final stage
FROM alpine:latest
# change apk source to ali
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk upgrade
# install: tzdata ca-certificates opencv-dev
RUN apk --no-cache add tzdata ca-certificates
# cpoy the app from builder
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=hominsu/crontab-worker Version=1.0
