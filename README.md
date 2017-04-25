# squirrel

## Deploying to Docker

Make sure you have [Go](https://golang.org/doc/install), [Glide](https://github.com/Masterminds/glide), and [go.rice](https://github.com/GeertJohan/go.rice) installed.

    $ glide install
    $ go generate ./boxes
    $ go build -o ./dist/squirrel-server
    $ docker build -t lavenderx/squirrel .
    $ docker run -it --rm -p 8081:80 lavenderx/squirrel
