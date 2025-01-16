package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 优化状态
type OptimizationStatus struct {
	IsEnabled      bool      // 是否启用优化
	CurrentProfile string    // 当前配置
	StartTime      time.Time // 开始时间
	MinecraftBW    int64     // MC带宽使用
	BackupPath     string    // 备份文件路径
}

var (
	currentStatus = OptimizationStatus{
		IsEnabled: false,
	}
)

// 仪表盘
func showDashboard() {
	fmt.Println("\n=== 系统状态 ===")
	fmt.Printf("优化状态: %s\n", getStatusString())
	if currentStatus.IsEnabled {
		if mc, err := getMinecraftProcess(); err == nil {
			fmt.Printf("游戏进程: %s (%s)\n", mc.Name, mc.Title)
			fmt.Printf("进程 ID: %s\n", mc.PID)
			fmt.Printf("内存使用: %s\n", formatBytes(mc.Memory))
			fmt.Printf("网络流量: %s/s\n", formatBytes(currentStatus.MinecraftBW))
		}
		fmt.Printf("当前配置: %s\n", currentStatus.CurrentProfile)
		fmt.Printf("运行时间: %s\n", time.Since(currentStatus.StartTime).Round(time.Second))
		if currentStatus.BackupPath != "" {
			fmt.Printf("备份位置: %s\n", currentStatus.BackupPath)
		}
	}
	fmt.Println("==============")
}

// 格式化
func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
}

// 流量监控
func monitorNetworkUsage() {
	go func() {
		retryCount := 0
		for {
			if currentStatus.IsEnabled {
				// 获取Minecraft进程流量
				mcBW := getMinecraftBandwidth()
				if mcBW == 0 && retryCount < 3 {
					// 如果获取失败，等待短暂时间后重试
					time.Sleep(500 * time.Millisecond)
					retryCount++
					continue
				}
				currentStatus.MinecraftBW = mcBW
				retryCount = 0

				// 检查并优化其他应用的网络使用
				optimizeNetworkPriority()
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// 获取MC使用量
func getMinecraftBandwidth() int64 {
	mc, err := getMinecraftProcess()
	if err != nil {
		return 0
	}

	// 使用进程ID获取网络使用情况
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		$process = Get-Process -Id %s
		$networkCounter = (Get-Counter "\Process($($process.Name))\IO Data Bytes/sec").CounterSamples.CookedValue
		Write-Output $networkCounter
	`, mc.PID))

	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	bw, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0
	}
	return bw
}

// 优化网络优先级
func optimizeNetworkPriority() {
	mc, err := getMinecraftProcess()
	if err != nil {
		return
	}

	// 2. 设置mc进程为最高优先级
	mcPriorityCmd := exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%s", mc.PID), "CALL", "setpriority", "realtime")
	if err := mcPriorityCmd.Run(); err != nil {
		fmt.Printf("警告: 设置Minecraft优先级失败: %v\n", err)
	}

	// 3. 获取其他高带宽使用的进程
	cmd := exec.Command("powershell", "-Command", `
		Get-NetTCPConnection | 
		Where-Object { $_.State -eq 'Established' } |
		Group-Object -Property OwningProcess | 
		Sort-Object -Property Count -Descending | 
		Select-Object -First 10 | 
		ForEach-Object { 
			$proc = Get-Process -Id $_.Name -ErrorAction SilentlyContinue
			if ($proc) {
				[PSCustomObject]@{
					PID = $_.Name
					Name = $proc.ProcessName
					Connections = $_.Count
				}
			}
		} |
		ConvertTo-Json
	`)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// 4. 解析进程信息
	type ProcessInfo struct {
		PID         string `json:"PID"`
		Name        string `json:"Name"`
		Connections int    `json:"Connections"`
	}

	var processes []ProcessInfo
	if err := json.Unmarshal(output, &processes); err != nil {
		return
	}

	// 5. 降低其他进程优先级
	for _, proc := range processes {
		if proc.PID == mc.PID {
			continue
		}

		if isSystemProcess(proc.Name) {
			continue
		}

		priorityCmd := exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%s", proc.PID), "CALL", "setpriority", "idle")
		priorityCmd.Run()
	}
}

func getStatusString() string {
	if currentStatus.IsEnabled {
		return "已启用"
	}
	return "未启用"
}

func enableOptimization(profile string) error {
	if currentStatus.IsEnabled {
		return fmt.Errorf("优化已经启用，请先关闭当前优化")
	}

	// 备份当前设置
	if err := backupCurrentSettings(); err != nil {
		return fmt.Errorf("备份设置失败: %v", err)
	}

	currentStatus.IsEnabled = true
	currentStatus.CurrentProfile = profile
	currentStatus.StartTime = time.Now()
	currentStatus.BackupPath = backupPath

	return nil
}

// 禁用优化
func disableOptimization() error {
	if !currentStatus.IsEnabled {
		return fmt.Errorf("优化未启用")
	}

	// 恢复网络设置
	regCommands := []RegCommand{
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpAckFrequency",
			Type:  "REG_DWORD",
			Value: "2",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpNoDelay",
			Type:  "REG_DWORD",
			Value: "0",
		},
	}

	for _, cmd := range regCommands {
		if err := executeRegCommand(cmd); err != nil {
			return fmt.Errorf("恢复设置失败: %v", err)
		}
	}

	// 只恢复自动调优级别
	cmd := exec.Command("netsh", "int", "tcp", "set", "global", "autotuninglevel=normal")
	if err := cmd.Run(); err != nil {
		fmt.Printf("警告: 恢复自动调优设置失败: %v\n", err)
	}

	currentStatus.IsEnabled = false
	currentStatus.CurrentProfile = ""
	currentStatus.MinecraftBW = 0

	return nil
}
