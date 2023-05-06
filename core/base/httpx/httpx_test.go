package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {
	res := &Request{
		Method:      http.MethodPost,
		Url:         "https://api.xxx.com",
		ContentType: ContentJson,
		Body:        FormatJsonData(map[string]string{"world": "1"}),
		Cookies: map[string]string{
			"token": "xxx",
		},
	}
	client := New(res)
	resp, err := client.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Content))
}
