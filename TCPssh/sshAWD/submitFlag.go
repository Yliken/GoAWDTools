package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var storeflagfile = "flaglist.txt" //flag存儲文件
var pattern = `flag{[^}]*}`        //flag格式

func sub(flag string) {
	// 定义curl命令
	s := fmt.Sprintf("{\"flag\": \"%s\"}", flag)
	cmd := exec.Command("curl", "-X", "POST", "", // flag提交地址
		"-H", "", //认证
		"-d", s)

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing curl command:", err)
		return
	}

	// 打印命令输出
	fmt.Println(string(output))
}

func main() {
	for {

		re := regexp.MustCompile(pattern)
		openfile, err := os.Open(storeflagfile)
		if err != nil {
			fmt.Println(err)
		}
		defer openfile.Close()

		scanner := bufio.NewScanner(openfile)
		for scanner.Scan() {
			line := scanner.Text()
			submatch := re.FindStringSubmatch(line)
			if submatch != nil {
				fmt.Println(submatch[0])
				sub(submatch[0])
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
