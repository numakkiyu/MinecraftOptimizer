package main

import (
	"fmt"
	"os/exec"
)

type NetworkOptimizer struct {
	Process     *MinecraftProcess
	PacketRules []PacketRule
	NetworkMode string
}

// PacketRule 数据包规则
type PacketRule struct {
	Protocol   string // TCP/UDP
	Direction  string // in/out
	Priority   int    // 优先级 0-7
	QoS        bool   // 服务质量
	BufferSize int    // 缓冲区大小
}

// 网络模式常量
const (
	MODE_PVP    = "pvp"    // PVP模式：最小延迟
	MODE_NORMAL = "normal" // 普通模式：平衡
	MODE_STABLE = "stable" // 稳定模式：低丢包
)

// NewNetworkOptimizer 创建新的网络优化
func NewNetworkOptimizer(process *MinecraftProcess) *NetworkOptimizer {
	return &NetworkOptimizer{
		Process:     process,
		NetworkMode: MODE_PVP,
		PacketRules: []PacketRule{
			{
				Protocol:   "TCP",
				Direction:  "out",
				Priority:   7,
				QoS:        true,
				BufferSize: 65535,
			},
			{
				Protocol:   "UDP",
				Direction:  "out",
				Priority:   7,
				QoS:        true,
				BufferSize: 65535,
			},
		},
	}
}

// ApplyNetworkOptimizations 应用网络优化
func (no *NetworkOptimizer) ApplyNetworkOptimizations() error {
	// 1. 设置QoS策略
	if err := no.setQoSPolicy(); err != nil {
		return fmt.Errorf("设置QoS策略失败: %v", err)
	}

	// 2. 设置数据包优先级
	if err := no.setPacketPriority(); err != nil {
		return fmt.Errorf("设置数据包优先级失败: %v", err)
	}

	// 3. 优化TCP参数
	if err := no.optimizeTCPParams(); err != nil {
		return fmt.Errorf("优化TCP参数失败: %v", err)
	}

	// 4. 设置网络缓冲区
	if err := no.setNetworkBuffers(); err != nil {
		return fmt.Errorf("设置网络缓冲区失败: %v", err)
	}

	return nil
}

// setQoSPolicy 设置QoS策略
func (no *NetworkOptimizer) setQoSPolicy() error {
	if no.Process.Path == "" {
		return fmt.Errorf("无法获取进程路径")
	}

	// 使用netsh设置QoS策略
	cmd := exec.Command("netsh", "qos", "add", "policy", "name=MinecraftQoS",
		fmt.Sprintf("apppath=\"%s\"", no.Process.Path),
		"profile=gamemode",
		"throttlerate=0")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("设置QoS策略失败: %v (%s)", err, string(output))
	}
	return nil
}

// setPacketPriority 设置数据包优先级
func (no *NetworkOptimizer) setPacketPriority() error {
	for _, rule := range no.PacketRules {
		cmd := exec.Command("netsh", "qos", "add", "rule",
			fmt.Sprintf("name=Minecraft%s%s", rule.Protocol, rule.Direction),
			fmt.Sprintf("protocol=%s", rule.Protocol),
			fmt.Sprintf("dir=%s", rule.Direction),
			fmt.Sprintf("priority=%d", rule.Priority))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// optimizeTCPParams 优化TCP参数
func (no *NetworkOptimizer) optimizeTCPParams() error {
	regCommands := []RegCommand{
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpAckFrequency",
			Type:  "REG_DWORD",
			Value: "1", // 立即确认
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpNoDelay",
			Type:  "REG_DWORD",
			Value: "1", // 禁用Nagle算法
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TCPInitialRTT",
			Type:  "REG_DWORD",
			Value: "2", // 最小RTT
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "DefaultTTL",
			Type:  "REG_DWORD",
			Value: "64", // 优化TTL
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "MaxUserPort",
			Type:  "REG_DWORD",
			Value: "65534", // 最大端口数
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpMaxDataRetransmissions",
			Type:  "REG_DWORD",
			Value: "3", // 减少重传次数
		},
	}

	for _, cmd := range regCommands {
		if err := executeRegCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

// setNetworkBuffers 设置网络缓冲区
func (no *NetworkOptimizer) setNetworkBuffers() error {
	netshCommands := []NetshCommand{
		{Category: "tcp", Setting: "global", Name: "autotuninglevel", Value: "restricted"},
		{Category: "tcp", Setting: "global", Name: "ecncapability", Value: "disabled"},
		{Category: "tcp", Setting: "global", Name: "timestamps", Value: "disabled"},
		{Category: "tcp", Setting: "global", Name: "rss", Value: "enabled"},
		{Category: "tcp", Setting: "global", Name: "initialRto", Value: "2000"},
	}

	for _, cmd := range netshCommands {
		if err := executeNetshCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}
