package main

import (
	"fmt"
	"os"
)

func main() {
	stat, err := os.Stat("ssh/Iplist.txt")
	if err != nil {
		// 如果发生错误，打印错误并退出程序
		fmt.Println("错误:", err)
		return
	}

	// 如果没有错误，安全地访问 stat
	fmt.Printf("文件名: %s\n", stat.Name())
}
