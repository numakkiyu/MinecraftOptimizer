package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// checkMinecraftRunning 检查MC是否运行
func checkMinecraftRunning() bool {
	_, err := getMinecraftProcess()
	return err == nil
}

// getMinecraftProcess 获取MC进程信息
func getMinecraftProcess() (*MinecraftProcess, error) {
	mcProcessNames := []string{
		"javaw.exe",
		"java.exe",
	}

	// PowerShell命令获取进程详细信息
	cmdText := `
		Get-Process | 
		Where-Object { $_.MainWindowTitle -like '*Minecraft*' -or $_.ProcessName -in @('javaw','java','Minecraft.Windows') } |
		Select-Object Id, ProcessName, MainWindowTitle, WorkingSet64, Path |
		ConvertTo-Json
	`

	cmd := exec.Command("powershell", "-Command", cmdText)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取进程信息失败: %v", err)
	}

	// 解析JSON输出
	type ProcessInfo struct {
		Id              int    `json:"Id"`
		ProcessName     string `json:"ProcessName"`
		MainWindowTitle string `json:"MainWindowTitle"`
		WorkingSet64    int64  `json:"WorkingSet64"`
		Path            string `json:"Path"`
	}

	var processes []ProcessInfo
	if err := json.Unmarshal(output, &processes); err != nil {
		return nil, fmt.Errorf("解析进程信息失败: %v", err)
	}

	// 查找Minecraft进程
	for _, proc := range processes {
		if containsAny(proc.ProcessName, mcProcessNames) ||
			strings.Contains(strings.ToLower(proc.MainWindowTitle), "minecraft") {
			return &MinecraftProcess{
				PID:    strconv.Itoa(proc.Id),
				Name:   proc.ProcessName,
				Title:  proc.MainWindowTitle,
				Memory: proc.WorkingSet64,
				Path:   proc.Path,
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到运行中的Minecraft进程")
}

// containsAny 检查字符串是否包含列表中的任何一个值
func containsAny(s string, list []string) bool {
	s = strings.ToLower(s)
	for _, item := range list {
		if strings.Contains(s, strings.ToLower(item)) {
			return true
		}
	}
	return false
}

// isSystemProcess 检查是否为系统关键进程
func isSystemProcess(name string) bool {
	systemProcesses := []string{
		"svchost", "system", "wininit", "winlogon", "services", "lsass",
		"csrss", "smss", "explorer", "dwm", "taskmgr",
	}
	name = strings.ToLower(name)
	for _, proc := range systemProcesses {
		if name == proc {
			return true
		}
	}
	return false
}
