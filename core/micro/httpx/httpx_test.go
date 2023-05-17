package httpx_test

import (
	"context"
	"fmt"
	"github.com/bang-go/crab/core/micro/httpx"
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {
	req := &httpx.Request{
		Method:      http.MethodPost,
		Url:         "https://api.xxx.com",
		ContentType: httpx.ContentJson,
		Body:        httpx.FormatJsonData(map[string]string{"world": "1"}),
		Cookies: map[string]string{
			"token": "xxx",
		},
	}
	client := httpx.New(httpx.Config{})
	resp, err := client.Send(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Content))
}
