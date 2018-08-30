FROM alpine:3.8
MAINTAINER Zongzhi Bai <dolphineor@gmail.com>

ENV LANG C.UTF-8

RUN echo "https://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories
RUN echo "https://mirrors.aliyun.com/alpine/v3.8/community/" >> /etc/apk/repositories

RUN apk update && apk --no-cache add bash ca-certificates tzdata vim && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

RUN mkdir -p /var/log/squirrel

ADD ./dist/squirrel /usr/local/squirrel/bin/squirrel

WORKDIR /usr/local/squirrel

CMD ["/bin/sh", "-c", "./bin/squirrel"]
