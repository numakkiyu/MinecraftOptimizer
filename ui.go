package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Bulletin 公告结构
type Bulletin struct {
	Content string `json:"bulletin"`
}

// getUserInput 获取用户输入
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// showMainMenu 显示主菜单
func showMainMenu() {
	clearScreen()
	showHeader()
	showAnnouncements()
	showDashboard()

	fmt.Println("\n=== 功能菜单 ===")
	fmt.Println("1. 基础优化")
	fmt.Println("2. 高级优化")
	fmt.Println("3. 应用设置")
	fmt.Println("4. 备份管理")
	fmt.Println("0. 退出程序")
	fmt.Println("============")
	fmt.Print("请输入数字选择功能: ")
}

// showBasicMenu 显示基础菜单
func showBasicMenu() {
	clearScreen()
	fmt.Println("=== 基础优化 ===")
	fmt.Println("1. 开启优化")
	fmt.Println("2. 关闭优化")
	fmt.Println("3. 切换模式")
	fmt.Println("0. 返回主菜单")
	fmt.Println("============")
}

// showModeMenu 显示模式选择菜单
func showModeMenu() {
	clearScreen()
	fmt.Println("=== 优化模式 ===")
	fmt.Println("1. 命中优化")
	fmt.Println("2. 击退优化")
	fmt.Println("3. 连接稳定性")
	fmt.Println("4. 叠刀优化")
	fmt.Println("0. 返回")
	fmt.Println("============")
}

// showAdvancedMenu 显示高级菜单
func showAdvancedMenu() {
	clearScreen()
	fmt.Println("=== 高级优化 ===")
	fmt.Println("1. 网络调优")
	fmt.Println("2. 开启叠刀")
	fmt.Println("3. 关闭叠刀")
	fmt.Println("4. 自定义配置")
	fmt.Println("0. 返回主菜单")
	fmt.Println("============")
}

// showBackupMenu 显示备份菜单
func showBackupMenu() {
	clearScreen()
	fmt.Println("=== 备份管理 ===")
	fmt.Println("1. 创建备份")
	fmt.Println("2. 恢复备份")
	fmt.Println("3. 查看备份")
	fmt.Println("0. 返回")
	fmt.Print("请选择: ")
}

// showTuningMenu 显示调优菜单
func showTuningMenu() {
	clearScreen()
	fmt.Println("=== 网络调优 ===")
	fmt.Println("1. 禁用自动调优")
	fmt.Println("2. 默认设置")
	fmt.Println("3. 受限模式")
	fmt.Println("4. 延迟模式")
	fmt.Println("5. 稳定模式")
	fmt.Println("0. 返回")
	fmt.Println("============")
}

// showHeader 显示头部信息
func showHeader() {
	fmt.Printf("=== Minecraft Optimizer 网络优化器 %s ===\n", Version)
	showCredits()
	fmt.Println("\n注意: 本程序属于GPL-3.0 license开源项目")
	fmt.Println("==================================")
}

// showAnnouncements 显示公告和广告
func showAnnouncements() {
	fmt.Println("\n=== 公告 ===")
	bulletin := fetchBulletin("https://api.tianbeigm.cn/tboz/bulletin.json")
	fmt.Println(bulletin)

	fmt.Println("\n=== 广告 ===")
	ad := fetchBulletin("https://api.tianbeigm.cn/tboz/ab.json")
	fmt.Println(ad)
}

// fetchBulletin 获取公告内容
func fetchBulletin(url string) string {
	client := &http.Client{
		Timeout: 5 * time.Second, // 添加超时设置
	}

	resp, err := client.Get(url)
	if err != nil {
		return "获取公告失败"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "读取公告失败"
	}

	var bulletin Bulletin
	if err := json.Unmarshal(body, &bulletin); err != nil {
		return "解析公告失败"
	}

	return bulletin.Content
}
