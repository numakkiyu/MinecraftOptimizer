package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	if err := loadStatus(); err != nil {
		fmt.Printf("警告: 加载状态失败: %v\n", err)
	}
}

func main() {
	// 检查管理员权限
	if !checkAdminPrivileges() {
		fmt.Println("请以管理员权限运行此程序！")
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		os.Exit(1)
	}

	// 检查MC是否运行
	if !checkMinecraftRunning() {
		fmt.Println("请先运行Minecraft游戏！")
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		os.Exit(1)
	}

	monitorNetworkUsage()

	// 程序退出时恢复注册表
	defer func() {
		if currentStatus.IsEnabled {
			if err := disableOptimization(); err != nil {
				fmt.Printf("警告: 恢复设置失败: %v\n", err)
			}
		}
	}()

	// 进入主菜单循环
	handler := &MainMenuHandler{}
	if err := handler.Handle(); err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

// 清屏函数
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}
