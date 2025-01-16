package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

// BackupInfo 备份信息
type BackupInfo struct {
	Time     time.Time        `json:"time"`
	Settings MinecraftProfile `json:"settings"`
}

// 备份当前网络设置
func backupCurrentSettings() error {
	backupDir := filepath.Join(os.TempDir(), "minecraft_optimizer")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	backupPath = filepath.Join(backupDir, fmt.Sprintf("registry_backup_%s.reg", time.Now().Format("20060102_150405")))

	// 备份TCP/IP参数
	cmd := exec.Command("reg", "export", `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, backupPath, "/y")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("备份注册表失败: %v, 输出: %s", err, output)
	}

	lastBackupTime = time.Now()
	if err := saveStatus(); err != nil {
		fmt.Printf("警告: 保存状态失败: %v\n", err)
	}

	return nil
}

// 获取当前网络设置
func getCurrentNetworkSettings() (MinecraftProfile, error) {
	// 读取当前注册表值
	settings := MinecraftProfile{
		Name:      "当前设置",
		RegParams: make([]RegCommand, 0),
		NetParams: make([]NetshCommand, 0),
	}

	// 读取关键注册表项
	regKeys := []string{
		"TcpAckFrequency",
		"DefaultTTL",
		"GlobalMaxTcpWindowSize",
		"TcpNoDelay",
		"MaxUserPort",
		"TcpMaxDataRetransmissions",
	}

	for _, key := range regKeys {
		value, err := readRegistryValue(key)
		if err != nil {
			continue // 跳过读取失败的项
		}
		settings.RegParams = append(settings.RegParams, RegCommand{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  key,
			Type:  "REG_DWORD",
			Value: value,
		})
	}

	return settings, nil
}

// 恢复默认网络设置
func restoreDefaultSettings() error {
	fmt.Println("正在恢复默认网络设置...")

	for _, cmd := range defaultNetworkSettings.RegParams {
		if err := executeRegCommand(cmd); err != nil {
			return fmt.Errorf("恢复默认注册表设置失败: %v", err)
		}
	}

	for _, cmd := range defaultNetworkSettings.NetParams {
		if err := executeNetshCommand(cmd); err != nil {
			return fmt.Errorf("恢复默认网络设置失败: %v", err)
		}
	}

	fmt.Println("已恢复默认网络设置")
	return nil
}

// 从备份恢复设置
func restoreFromBackup() error {
	backups, err := listBackupFiles()
	if err != nil {
		return err
	}

	if len(backups) == 0 {
		return fmt.Errorf("没有找到备份文件")
	}

	latestBackup := backups[len(backups)-1]
	fmt.Printf("正在从备份文件恢复: %s\n", filepath.Base(latestBackup))

	cmd := exec.Command("reg", "import", latestBackup)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("恢复注册表失败: %v, 输出: %s", err, output)
	}

	return nil
}

// 列出所有备份
func listBackups() error {
	backups, err := listBackupFiles()
	if err != nil {
		return err
	}

	if len(backups) == 0 {
		fmt.Println("没有找到备份文件")
		return nil
	}

	fmt.Println("\n可用的备份文件:")
	for i, backup := range backups {
		info, err := os.Stat(backup)
		if err != nil {
			continue
		}
		fmt.Printf("%d. %s (创建于: %s)\n", i+1, filepath.Base(backup), info.ModTime().Format("2006-01-02 15:04:05"))
	}

	return nil
}

// 获取所有备份文件列表
func listBackupFiles() ([]string, error) {
	backupDir := filepath.Join(os.TempDir(), "minecraft_optimizer")
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".reg" {
			backups = append(backups, filepath.Join(backupDir, file.Name()))
		}
	}

	// 按修改时间排序
	sort.Slice(backups, func(i, j int) bool {
		info1, _ := os.Stat(backups[i])
		info2, _ := os.Stat(backups[j])
		return info1.ModTime().Before(info2.ModTime())
	})

	return backups, nil
}
