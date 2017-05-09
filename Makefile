# https://seisman.github.io/how-to-write-makefile/index.html
BINARY_FILE=./dist/squirrel-server
IMAGE=lavenderx/squirrel-caddy

build:
	cd ./ui && npm run gulp build-prod && cd -
	go generate ./app
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_FILE)
	docker build -f Dockerfile.caddy -t $(IMAGE) .

clean:
	if [ -f $(BINARY_FILE) ]; then rm -f $(BINARY_FILE); fi
	go clean

exec: clean build
