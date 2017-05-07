#!/bin/bash

# Installing Golang dependencies
glide install

# Installing ui dependencies
cd ./ui && npm install --registry=https://registry.npm.taobao.org

# Building ui resources
npm run gulp build-prod && cd ..

# Embedding static files
go generate ./app

# Building an executable binary file
# https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./dist/squirrel-server

if [ $1 == "mini" ]; then
    # Building a minimal image, about 20MB.
    docker build -f Dockerfile.alpine -t lavenderx/squirrel-alpine .
    # Running a temporary container to start application
    docker run -it --rm -p 8081:7000 lavenderx/squirrel-alpine
else
    # Building a normal image
    docker build -t lavenderx/squirrel .
    docker run -it --rm -p 8081:80 lavenderx/squirrel
fi
