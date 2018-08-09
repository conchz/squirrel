FROM alpine:3.5
MAINTAINER Zongzhi Bai <dolphineor@gmail.com>

ENV PYTHON_VERSION=2.7.14-r0
ENV PY2_PIP_VERSION=9.0.0-r1
ENV SUPERVISOR_VERSION=3.3.1

# Setup TimeZone
RUN apk update && apk --no-cache add bash ca-certificates tzdata vim && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# Install Supervisor
RUN apk --no-cache add -u python=$PYTHON_VERSION py2-pip=$PY2_PIP_VERSION
RUN pip install supervisor==$SUPERVISOR_VERSION

# Create log directory
RUN mkdir -p /var/log/supervisor && \
    mkdir -p /var/log/squirrel

ADD ./dist/squirrel-server /usr/local/bin/
ADD ./supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 7000

ENTRYPOINT ["supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
