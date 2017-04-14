# squirrel

## Deploying to Docker

Make sure you have [Go](https://golang.org/doc/install), [Glide](https://github.com/Masterminds/glide), and the [go.rice](https://github.com/GeertJohan/go.rice) installed.

    $ glide install
    $ rice embed-go
    $ go build
    $ docker build -t lavenderx/squirrel .
    $ docker run -it --rm -p 8081:80 lavenderx/squirrel
