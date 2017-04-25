#!/bin/bash

# Installing project dependencies
glide install

# Embedding static files
go generate ./boxes

# Building an executable binary file
go build -o ./dist/squirrel-server

# Building a docker image
docker build -t lavenderx/squirrel .

# Running a temporary docker container to start application
docker run -it --rm -p 8081:80 lavenderx/squirrel