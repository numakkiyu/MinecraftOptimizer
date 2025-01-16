package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// 优化本地网络延迟
func optimizeLocalDelay() error {
	if err := optimizeNetworkAdapter(); err != nil {
		return err
	}

	// 2. 优化TCP参数
	regCommands := []RegCommand{
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpAckFrequency",
			Type:  "REG_DWORD",
			Value: "1",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpNoDelay",
			Type:  "REG_DWORD",
			Value: "1",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpDelAckTicks",
			Type:  "REG_DWORD",
			Value: "0",
		},
	}

	for _, cmd := range regCommands {
		if err := executeRegCommand(cmd); err != nil {
			return err
		}
	}

	// 3. 优化网络服务
	netshCommands := []string{
		"int tcp set global autotuninglevel=restricted",
		"int tcp set supplemental congestionprovider=ctcp",
		"int tcp set global ecncapability=disabled",
		"int tcp set global timestamps=disabled",
		"int tcp set global rss=enabled",
		"int tcp set global initialRto=2000",
	}

	for _, cmdStr := range netshCommands {
		cmd := exec.Command("netsh", strings.Split(cmdStr, " ")...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("执行网络命令失败: %v, 输出: %s", err, output)
		}
	}

	return nil
}

// 优化网卡设置
func optimizeNetworkAdapter() error {
	// 设置网卡优先级
	cmd := exec.Command("wmic", "nicconfig", "where", "IPEnabled=TRUE", "call", "SetPowerManagement", "FALSE")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置网卡电源管理失败: %v", err)
	}
	return nil
}

// 优化网络缓冲区
func optimizeNetworkBuffer() error {
	regCommands := []RegCommand{
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "GlobalMaxTcpWindowSize",
			Type:  "REG_DWORD",
			Value: "65535",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpWindowSize",
			Type:  "REG_DWORD",
			Value: "65535",
		},
	}

	for _, cmd := range regCommands {
		if err := executeRegCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}
