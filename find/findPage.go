package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	ipfile := "Iplist.txt"
	open, err := os.Open(ipfile)
	if err != nil {
		fmt.Println("open ipfile error", err)
	}
	defer open.Close()
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()
		resp, err := http.Get(fmt.Sprintf("%s/", line))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Println(fmt.Sprintf("%s/", line), "存在")
		}
	}
}
