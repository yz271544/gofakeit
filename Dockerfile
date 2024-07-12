# build stage
FROM golang:1.22.4-alpine as builder

MAINTAINER Hu Lyndon <huzhengyang@gridsum.com>

WORKDIR /code
COPY . .
RUN go build -v -o /usr/local/bin/gofakeitserver -ldflags "-w -s" \
./cmd/gofakeitserver/main.go

# final stage
FROM alpine:3.20
MAINTAINER Hu Lyndon <huzhengyang@gridsum.com>

COPY --from=builder usr/local/bin/gofakeitserver /usr/local/bin/gofakeitserver
EXPOSE 8080
CMD ["gofakeitserver"]
