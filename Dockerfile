# build stage
FROM golang:1.24-alpine as builder

MAINTAINER Hu Lyndon <huzhengyang@gridsum.com>

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn|https://mirrors.aliyun.com/goproxy/|https://goproxy.io|direct \
    GOPRIVATE=*.gitlab.com,*.gitee.com,*.github.com \
    CGO_ENABLED=0

WORKDIR /code
COPY . .
RUN go mod tidy
RUN go build -v -o /usr/local/bin/gofakeitserver -ldflags "-w -s" cmd/gofakeitserver/main.go

# final stage
FROM alpine:3.22.0
MAINTAINER Hu Lyndon <huzhengyang@gridsum.com>

COPY --from=builder usr/local/bin/gofakeitserver /usr/local/bin/gofakeitserver
EXPOSE 8080
CMD ["gofakeitserver"]
