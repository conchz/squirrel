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
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./dist/squirrel-server

docker build -t lavenderx/squirrel .
docker run --name squirrel-demo -it --rm -p 8081:7000 lavenderx/squirrel
