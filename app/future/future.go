package future

import (
	"io/ioutil"
	"net/http"
)

// Exec SimpleFuture to get the response. It's an implementation of Golang Future/Promise.
// reference: https://yushuangqi.com/blog/2017/golangli-de-future_promise.html
//
//   Params:
//     url: request url
//   Returns:
//     response body
//     error, if error occurred
func SimpleFuture(url string) func() ([]byte, error) {
	var body []byte
	var err error

	c := make(chan struct{}, 1)
	go func() {
		defer close(c)

		var res *http.Response
		res, err = http.Get(url)
		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
	}()

	return func() ([]byte, error) {
		<-c
		return body, err
	}
}

// Future boilerplate method
func Future(f func() (interface{}, error)) func() (interface{}, error) {
	var result interface{}
	var err error

	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()

	return func() (interface{}, error) {
		<-c
		return result, err
	}
}
