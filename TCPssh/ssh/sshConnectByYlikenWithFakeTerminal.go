package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	//本地客户端配置
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("312909"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//与远程主机建立连接
	dial, err := ssh.Dial("tcp", "192.168.22.100:22", config)
	if err != nil {
		fmt.Println("与远程主机连接失败,请检查网络ip或者端口是否输入错误!")
	}
	defer dial.Close()

	//与远程主机建立会话
	session, err := dial.NewSession()
	if err != nil {
		fmt.Println("建立会话失败")
	}
	defer session.Close()

	//新建伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Println(err)
	}

	//绑定标准输入与标准输出
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	session.Shell()

	session.Wait()

}
