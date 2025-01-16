package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	// 备份相关
	backupPath     string        // 备份文件路径
	lastBackupTime time.Time     // 最后备份时间
	backupInterval time.Duration // 备份间隔
)

func init() {
	tempDir := os.TempDir()

	backupPath = filepath.Join(tempDir, "minecraft_optimizer", "backups")

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		fmt.Printf("警告: 创建备份目录失败: %v\n", err)
	}

	backupInterval = 24 * time.Hour

	if err := loadStatus(); err != nil {
		fmt.Printf("警告: 加载状态失败: %v\n", err)
	}
}

// saveStatus 保存当前状态到文件
func saveStatus() error {
	statusDir := filepath.Join(backupPath, "status")
	if err := os.MkdirAll(statusDir, 0755); err != nil {
		return fmt.Errorf("创建状态目录失败: %v", err)
	}

	data, err := json.Marshal(currentStatus)
	if err != nil {
		return err
	}

	statusFile := filepath.Join(statusDir, "status.json")
	return os.WriteFile(statusFile, data, 0644)
}

// executeRegCommand 执行注册表命令
func executeRegCommand(cmd RegCommand) error {
	args := []string{
		"add",
		cmd.Path,
		"/v", cmd.Name,
		"/t", cmd.Type,
		"/d", cmd.Value,
		"/f",
	}
	command := exec.Command("reg", args...)
	if err := command.Run(); err != nil {
		return fmt.Errorf("执行注册表命令失败: %v", err)
	}
	return nil
}

// executeNetshCommand 执行网络命令
func executeNetshCommand(cmd NetshCommand) error {
	args := []string{
		"int",
		cmd.Category,
		"set",
		cmd.Setting,
		fmt.Sprintf("%s=%s", cmd.Name, cmd.Value),
	}
	command := exec.Command("netsh", args...)
	if err := command.Run(); err != nil {
		return fmt.Errorf("执行网络命令失败: %v", err)
	}
	return nil
}

// loadStatus 从文件加载状态
func loadStatus() error {
	statusFile := filepath.Join(backupPath, "status", "status.json")
	data, err := os.ReadFile(statusFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &currentStatus)
}
