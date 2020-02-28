package request

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

// HTTPTransport lib make request
type HTTPTransport struct {
	StatusCode int
	BodyByte   []byte
	Err        error
	*fasthttp.Response
}

// NewHTTPTransport new instance
func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{}
}

// GET method get
func (ts *HTTPTransport) GET(url string) *HTTPTransport {
	ts.StatusCode, ts.BodyByte, ts.Err = fasthttp.Get(nil, url)
	return ts
}

// ToJSON body []byte to struct
func (ts *HTTPTransport) ToJSON() (interface{}, error) {
	var jsonCtx interface{}
	err := json.Unmarshal(ts.BodyByte, &jsonCtx)
	if err != nil {
		log.Println(`[utils.HTTPTransport] error:`, err)
	}
	return jsonCtx, err
}
