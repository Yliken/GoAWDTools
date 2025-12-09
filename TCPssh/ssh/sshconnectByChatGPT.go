package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
)

var user = "root"
var password = "ChengSiHannb666"

func main() {
	// SSH 配置
	config := &ssh.ClientConfig{
		User: user, // SSH 用户名
		Auth: []ssh.AuthMethod{
			ssh.Password(password), // SSH 密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 跳过主机密钥验证（生产环境请使用更安全的方式）
	}

	// 连接到远程服务器
	client, err := ssh.Dial("tcp", "192.168.22.100:22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput("ls -l /")
	if err != nil {
		log.Fatalf("Failed to run command: %s", err)
	}

	// 打印输出
	fmt.Println(string(output))
}
