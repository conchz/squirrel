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
go build -o ./dist/squirrel-server

# Building a docker image
docker build -t lavenderx/squirrel .

# Running a temporary docker container to start application
docker run -it --rm -p 8081:80 lavenderx/squirrel
