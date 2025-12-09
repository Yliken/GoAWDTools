# AWD 工具集 - Go 语言实现

本仓库包含一套使用 **Go 语言** (Golang) 编写的工具集，旨在协助 **攻防对抗 (AWD)** 竞赛。这些工具专注于 AWD 场景中的常见任务，例如服务发现、通过 SSH 进行远程命令执行以及 Web Shell 管理。

## 🌟 功能特性

该工具集按模块组织，每个模块都针对 AWD 竞赛中的特定需求：

| 模块             | 主要功能                  | 关键组件                                                                                                      |
| :--------------- | :------------------------ |:----------------------------------------------------------------------------------------------------------|
| **`find`**       | **服务发现**              | 一个脚本 (`findWensite.go`)，用于扫描目标 IP 上的端口范围，以识别正在运行的 Web 服务（通过状态码 200 判断）。                                   |
| **`TCPssh`**     | **通过 SSH 进行远程管理** | 一系列通过 SSH 管理目标主机的脚本，包括：                                                                                   |
|                  |                           | - **密码管理**：例如 `sshChangeMe.go` 等脚本，用于自动化修改远程主机的密码。`changeallpasswd`则尝试修改目标网络所有主机的密码                       |
|                  |                           | - **Flag 提交**：用于读取和提交 Flag 的脚本（`submitFlag.go`, `readflag.go`）。                                           |
|                  |                           | - **命令执行**：用于在远程主机上执行命令的脚本，具有多种连接方式（例如 `sshconnectByChatGPT.go`, `sshConnectByYlikenScript.go`）。          |
| **`writeHorse`** | **Web Shell 管理**        | 一个脚本 (`byPOST.go`)，用于通过 HTTP POST 请求与多个目标上的已知 Web Shell 或命令执行漏洞（例如 `admin_ccmd` 参数）进行交互。然后进行批量写webshell操作 |

## 🛠️ 安装指南

本项目使用 Go 语言编写，需要 Go 环境才能构建和运行。

### 先决条件

*   **Go (Golang)**: 推荐使用 1.20 或更高版本。

### 构建工具

1. **克隆仓库**（假设您会将此项目托管在 GitHub 等平台上）：

   ```bash
   git clone [YOUR_REPOSITORY_URL]
   cd AWD
   ```

2. **构建可执行文件**：
   由于项目结构中包含多个 `main` 包，您需要单独构建每个工具。

   ```bash
   # 示例：构建服务发现工具
   go build -o findWebsite ./find/findWensite.go
   
   # 示例：构建 Web Shell 交互工具
   go build -o writeHorse ./writeHorse/byPOST.go
   
   # 示例：构建 SSH 密码修改工具
   go build -o sshChangeMe ./TCPssh/sshAWD/sshChangeMe.go
   ```

## 🚀 使用示例

### 1. 服务发现 (`findWensite.go`)

此工具用于快速扫描目标机器上运行的 Web 服务。**注意**：在运行之前，您需要修改 `findWensite.go` 文件中硬编码的 IP 地址和端口范围。

```bash
./findWebsite
```

### 2. Web Shell 交互 (`byPOST.go`)

此工具使用特定的 POST 参数 (`admin_ccmd`) 向目标 Web 服务列表发送命令。**注意**：您需要修改 `byPOST.go` 文件中的 `target` 数组和 `payload`。

```bash
./writeHorse
```

### 3. SSH 密码修改 (`sshChangeMe.go`)

此工具通过 SSH 自动化修改远程主机上用户的密码。**重要**：在执行之前，您必须更新 `sshChangeMe.go` 中的以下变量：

*   `oldMypassword`: 当前密码。
*   `newMypassword`: 要设置的新密码。
*   `Myserver`: 目标 IP 和端口（例如 `"192.168.22.101:22"`）。
*   `User` in `sshConfig`: 目标用户名（例如 `"Yliken"`）。

```bash
./sshChangeMe
```

## ⚠️ 免责声明

这些工具仅供在受控、合法和道德的环境中使用，例如经过授权的夺旗赛 (CTF) 或攻防对抗赛 (AWD)。作者对因滥用这些工具而造成的任何损害不承担任何责任。**请自行承担风险，并且仅在您拥有明确测试权限的系统上使用。**
