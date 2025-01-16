package main

import (
	"fmt"
	"runtime"
)

// 版本信息
var (
	Version   = "v0.2.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = runtime.Version()
)

// 开发者信息
var (
	Developer = "北海的佰川"
	Credits   = []string{
		"Exambia",
		"Dian_Lu",
	}
	Repository = "https://github.com/numakkiyu/MinecraftOptimizer"
	Website    = "https://me.tianbeigm.cn"
	License    = "GPL-3.0"
)

// showVersion 显示版本信息
func showVersion() {
	fmt.Printf("Minecraft Optimizer 网络优化器 %s\n", Version)
	fmt.Printf("构建时间: %s\n", BuildTime)
	fmt.Printf("Git 提交: %s\n", GitCommit)
	fmt.Printf("Go 版本: %s\n", GoVersion)
}

// showCredits 显示开发者信息
func showCredits() {
	fmt.Printf("开发者: %s\n", Developer)
	fmt.Print("鸣谢: ")
	for i, credit := range Credits {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(credit)
	}
	fmt.Println()
	fmt.Printf("项目地址: %s\n", Repository)
	fmt.Printf("官网: %s\n", Website)
	fmt.Printf("许可证: %s\n", License)
}
