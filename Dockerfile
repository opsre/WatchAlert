FROM golang:1.23.0-alpine3.20 AS build

ARG VERSION

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /root

COPY . /root

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx tzdata

RUN CGO_ENABLED=0 go build --ldflags="-X main.Version=${VERSION}" -o w8t . \
    && chmod +x w8t

FROM alpine:3.19

COPY --from=build /root/w8t /app/w8t

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx tzdata

WORKDIR /app

ENTRYPOINT ["/app/w8t"]
