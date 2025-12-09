package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	Qaddr := "10.50.246."
	for i := 1; i <= 254; i++ {
		status := CheckWebsiteStatus(fmt.Sprintf("http://%s%d/.ccc.php", Qaddr, i))
		if status == true {
			fmt.Println(fmt.Sprintf("http://%s%d/.ccc.php", Qaddr, i))
		}
	}
}
func CheckWebsiteStatus(url string) bool {
	// 创建一个带超时的 HTTP 客户端
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
