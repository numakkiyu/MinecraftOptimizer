package main

// MinecraftProfile 定义MC网络配置
type MinecraftProfile struct {
	Name        string
	Description string
	RegParams   []RegCommand
	NetParams   []NetshCommand
	BufferSize  int
	Priority    string
}

// MC网络配置集
var minecraftProfiles = map[string]MinecraftProfile{
	"pvp": {
		Name:        "PVP优化",
		Description: "优化PVP网络延迟和击退",
		RegParams: []RegCommand{
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpAckFrequency",
				Type:  "REG_DWORD",
				Value: "1", // 立即确认，减少延迟
			},
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpNoDelay",
				Type:  "REG_DWORD",
				Value: "1", // 禁用Nagle算法
			},
			// 发包优化
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "DefaultTTL",
				Type:  "REG_DWORD",
				Value: "64", // 优化TTL，减少路由跳数
			},
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "MaxUserPort",
				Type:  "REG_DWORD",
				Value: "65534", // 增加可用端口数
			},
			// 网络缓冲区优化
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "GlobalMaxTcpWindowSize",
				Type:  "REG_DWORD",
				Value: "65535", // 最大TCP窗口大小
			},
			// 重传和超时设置
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpMaxDataRetransmissions",
				Type:  "REG_DWORD",
				Value: "3", // 减少重传次数，加快响应
			},
		},
		NetParams: []NetshCommand{
			{Category: "tcp", Setting: "autotuninglevel", Value: "restricted"},
			{Category: "tcp", Setting: "ecncapability", Value: "disabled"},
			{Category: "tcp", Setting: "timestamps", Value: "disabled"},
			{Category: "tcp", Setting: "initialRto", Value: "2000"},
			{Category: "tcp", Setting: "congestionprovider", Value: "ctcp"},
		},
		BufferSize: 65535,
		Priority:   "high",
	},
	"nokb": {
		Name:        "无击退优化",
		Description: "最小化击退效果",
		RegParams: []RegCommand{
			// 延迟确认设置
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpAckFrequency",
				Type:  "REG_DWORD",
				Value: "3", // 延迟确认
			},
			// 发包控制
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpMaxDataRetransmissions",
				Type:  "REG_DWORD",
				Value: "2", // 最小重传
			},
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "DefaultTTL",
				Type:  "REG_DWORD",
				Value: "48", // 较低TTL
			},
			// 缓冲区设置
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "GlobalMaxTcpWindowSize",
				Type:  "REG_DWORD",
				Value: "32768", // 较小的TCP窗口
			},
		},
		NetParams: []NetshCommand{
			{Category: "tcp", Setting: "autotuninglevel", Value: "disabled"},
			{Category: "tcp", Setting: "ecncapability", Value: "disabled"},
			{Category: "tcp", Setting: "congestionprovider", Value: "none"},
		},
		BufferSize: 32768,
		Priority:   "low",
	},
}

// 默认网络设置
var defaultNetworkSettings = MinecraftProfile{
	Name:        "默认设置",
	Description: "Windows默认网络配置",
	RegParams: []RegCommand{
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpAckFrequency",
			Type:  "REG_DWORD",
			Value: "2",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "DefaultTTL",
			Type:  "REG_DWORD",
			Value: "128",
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "GlobalMaxTcpWindowSize",
			Type:  "REG_DWORD",
			Value: "16384",
		},
	},
	NetParams: []NetshCommand{
		{Category: "tcp", Setting: "autotuninglevel", Value: "normal"},
		{Category: "tcp", Setting: "congestionprovider", Value: "default"},
		{Category: "tcp", Setting: "ecncapability", Value: "enabled"},
	},
	BufferSize: 16384,
	Priority:   "normal",
}
