package httpx

import (
	"github.com/bang-go/opt"
	"net/http"
	"time"
)

const (
	ContentRaw     = "Raw"            //原始请求
	ContentForm    = "Form"           //Form请求
	ContentJson    = "Json"           //Json请求
	DefaultTimeout = 30 * time.Second //默认请求时间

)

type Client interface {
	Send(opts ...opt.Option[requestOptions]) (resp *Response, err error)
}

type ClientWrapper struct {
	request *Request
}

func New(r *Request) Client {
	return &ClientWrapper{
		request: r,
	}
}
func (c *ClientWrapper) Send(opts ...opt.Option[requestOptions]) (resp *Response, err error) {
	options := &requestOptions{}
	opt.Each(options, opts...)
	httpUrl, err := c.request.getUrl()
	if err != nil {
		return
	}
	method, err := c.request.getMethod()
	if err != nil {
		return
	}
	reqBody := c.request.getBody()
	var httpReq *http.Request
	var httpRes *http.Response
	if httpReq, err = http.NewRequest(method, httpUrl, reqBody); err != nil { //新建http请求
		return
	}
	c.request.setHeaders(httpReq) //init headers
	//basic auth
	if options.baseAuth != nil {
		httpReq.SetBasicAuth(options.baseAuth.Username, options.baseAuth.Password)
	}
	c.request.setCookie(httpReq) ////init cookie
	httpClient := &http.Client{}
	// set timeout
	httpClient.Timeout = c.request.Timeout
	if httpClient.Timeout == 0 {
		httpClient.Timeout = DefaultTimeout
	}

	startTime := time.Now()
	if httpRes, err = httpClient.Do(httpReq); err != nil {
		return
	}
	defer httpRes.Body.Close()
	endTime := time.Now()
	elapsed := endTime.Sub(startTime).Seconds()
	resp = c.request.packResponse(httpRes, elapsed)
	return
}
