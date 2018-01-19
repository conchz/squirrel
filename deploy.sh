#!/bin/bash

# Installing Golang dependencies
dep ensure

# Installing ui dependencies
cd ./ui && npm install --registry=https://registry.npm.taobao.org

# Building ui resources
npm run gulp build-prod && cd ..

# Embedding static files
go generate ./app

# Building an executable binary file
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./dist/squirrel-server

# Building docker image
docker build -t lavenderx/squirrel .

containerName=squirrel-demo
isRunning=$(docker inspect --format="{{ .State.Running }}" ${containerName} 2> /dev/null)

if [ $? -eq 1 ];  then
    printf "No ${containerName} is running\n"
else
    # Checking container running status
    if [ "${isRunning}" = "true" ];  then
        docker stop ${containerName}
    fi

    docker rm ${containerName}
fi

docker run --name ${containerName} --net host -d lavenderx/squirrel
