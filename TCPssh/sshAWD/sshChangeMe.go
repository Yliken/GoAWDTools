package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

var oldMypassword = ""
var newMypassword = ""
var Myserver = "192.168.22.101:22" //要改 加 端口

func changeMyPassword(oldpassword, newpassword, server string) bool { // SSH 连接配置
	sshConfig := &ssh.ClientConfig{
		User:            "Yliken",                                    //要改                              // 填写 SSH 用户名
		Auth:            []ssh.AuthMethod{ssh.Password(oldpassword)}, // 填写 SSH 密码
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),                 // 忽略主机密钥验证
		Timeout:         30 * time.Second,
	}

	// 创建与服务器的连接
	client, err := ssh.Dial("tcp", server, sshConfig)
	if err != nil {
		fmt.Println(server, "连接失败")
		return false
	}
	defer client.Close()

	// 开始一个新的 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
		return false
	}
	defer session.Close()

	// 使用 passwd 命令修改密码
	cmd := fmt.Sprintf("echo -e '%s\\n%s\\n%s' | passwd", oldpassword, newpassword, newpassword)

	// 执行命令echo -e 'yliken\n123456\n123456' | passwd
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatalf("Failed to change password: %v", err)
		return false
	}

	// 输出结果
	fmt.Println(string(output))
	return true

}

func main() {
	changeMyPassword(oldMypassword, newMypassword, Myserver)
}
