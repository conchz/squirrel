### App deployment script to create a new LXC Container via Docker
###
### Docker: https://www.docker.com

FROM nginx:1.10.3
MAINTAINER Zongzhi Bai <dolphineor@gmail.com>

# Tell debconf to run in non-interactive mode
ENV DEBIAN_FRONTEND noninteractive

# Update & Install System Dependencies
RUN apt-get update && \
    apt-get -y install build-essential curl vim python-pip python-setuptools

# Install & Verify Go
ENV GOLANG_VERSION 1.8.1
WORKDIR /root
RUN mkdir -p /root/go/bin
RUN curl -qO https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz \
    && tar -xzf go$GOLANG_VERSION.linux-amd64.tar.gz -C /usr/local \
    && rm -f go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOROOT /usr/local/go
ENV GOPATH /root/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
RUN go env

# Install Supervisor
RUN mkdir ~/.pip
RUN echo "[global]\nindex-url=http://mirrors.aliyun.com/pypi/simple/\n[install]\ntrusted-host=mirrors.aliyun.com" > ~/.pip/pip.conf
RUN pip install supervisor
RUN pip install supervisor-stdout

# Set Time Zone
RUN echo "Asia/Shanghai" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

# Stage App
ADD ./dist/squirrel-server $GOPATH/bin

# Add Nginx frontend host
ADD ./docker/nginx-echo.vhost /etc/nginx/conf.d/default.conf

# Setup Nginx
RUN sed -i -e"s/keepalive_timeout\s*65/keepalive_timeout 3/" /etc/nginx/nginx.conf
RUN sed -i -e"s/keepalive_timeout 2/keepalive_timeout 3;\n\tclient_max_body_size 100m/" /etc/nginx/nginx.conf
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

# Setup Supervisord
ADD ./docker/supervisord.conf /etc/supervisord.conf

# Set start script permissions
ADD ./docker/startup.sh /startup.sh
RUN chmod 755 /startup.sh

# Expose Web Frontend (Nginx) port only
EXPOSE 80

# Start required services when docker is instantiated
ENTRYPOINT ["/startup.sh"]
