package main

import (
	"go_shell/config"
	"go_shell/ssh"
)

// main 本地连接Linux服务，执行shell命令
func main() {
	mConfig := config.LoadConfig("./config.yaml")
	go func(ip, port, user, pwd string) {
		ssh.SSHOpen(ip+":"+port, user, pwd)
	}(mConfig.Server.Ip, mConfig.Server.Port, mConfig.Server.User, mConfig.Server.Pwd)

	ssh.Execute("cd /home")

	// 保持不退出
	select {}
}
