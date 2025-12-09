package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	payload := "hello"
	values := url.Values{}
	values.Add("admin_ccmd", payload)
	for i := 0; i < 65535; i++ {
		req, err := http.Post(fmt.Sprintf("http://1.94.188.98:%d", i), "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
		if err != nil {
			continue
		}
		if req.StatusCode == 200 {
			fmt.Println(fmt.Sprintf("http://1.94.188.98:%d", i))
		}
	}
}
