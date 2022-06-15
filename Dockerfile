FROM golang:1.12.5-alpine as builder
ARG VERSION
RUN mkdir -p /go/src/github.com/TimothyYe/bing-wallpaper
WORKDIR /go/src/github.com/TimothyYe/bing-wallpaper
RUN cd /go/src/github.com/TimothyYe/bing-wallpaper
COPY . .
RUN apk --no-cache add git build-base make gcc libtool musl-dev ca-certificates dumb-init \
    && GO111MODULE=on go build -o ./bw/bw ./bw/main.go


FROM alpine
LABEL maintainer="bw"
RUN apk --no-cache add ca-certificates tzdata sqlite \
	&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& echo "Asia/Shanghai" >  /etc/timezone \
	&& apk del tzdata
# See https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN RUN ARCH=`uname -m` && mkdir /lib64 && ln -s "/lib/libc.musl-$ARCH.so.1" "/lib64/ld-linux-$ARCH.so.2"

RUN mkdir /bw
WORKDIR /bw
COPY --from=builder /go/src/github.com/TimothyYe/bing-wallpaper/bw/bw /bw/bw

EXPOSE 9000
ENTRYPOINT ["/bw/bw", "run"]
