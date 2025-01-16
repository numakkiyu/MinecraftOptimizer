package main

var (
	pvpProfiles = map[string]NetworkProfile{
		"hit_reg": {
			Name:        "命中优化",
			Description: "优化击打判定和命中反馈",
			RegParams: []RegCommand{
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
			},
			NetParams: []NetshCommand{
				{Category: "tcp", Setting: "global", Name: "autotuninglevel", Value: "restricted"},
				{Category: "tcp", Setting: "global", Name: "ecncapability", Value: "disabled"},
			},
			BufferSize: 65535,
			Priority:   "high",
		},
		"knockback": {
			Name:        "击退优化",
			Description: "优化击退效果和同步",
			RegParams: []RegCommand{
				{
					Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
					Name:  "TCPInitialRTT",
					Type:  "REG_DWORD",
					Value: "2",
				},
				{
					Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
					Name:  "DefaultTTL",
					Type:  "REG_DWORD",
					Value: "64",
				},
			},
			NetParams: []NetshCommand{
				{Category: "tcp", Setting: "global", Name: "timestamps", Value: "disabled"},
				{Category: "tcp", Setting: "global", Name: "rss", Value: "enabled"},
			},
			BufferSize: 131072,
			Priority:   "realtime",
		},
	}
)

// 网络调优设置
var tuningProfiles = map[string]NetworkProfile{
	"stability": {
		Name:        "连接稳定性",
		Description: "优化连接稳定性，减少断连",
		RegParams: []RegCommand{
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "TcpMaxDataRetransmissions",
				Type:  "REG_DWORD",
				Value: "5", // 增加重传次数
			},
			{
				Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
				Name:  "KeepAliveTime",
				Type:  "REG_DWORD",
				Value: "300000", // 5分钟保活
			},
		},
		NetParams: []NetshCommand{
			{Category: "tcp", Setting: "global", Name: "autotuninglevel", Value: "normal"},
			{Category: "tcp", Setting: "global", Name: "ecncapability", Value: "enabled"},
		},
		BufferSize: 131072,
		Priority:   "normal",
	},
}

// 叠刀设置 - 通过控制发包速度来实现快速连击
var stackedHitSettings = map[string][]RegCommand{
	"on": {
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
			Name:  "TcpDelAckTicks",
			Type:  "REG_DWORD",
			Value: "0", // 最小延迟
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TCPInitialRTT",
			Type:  "REG_DWORD",
			Value: "2", // 最小RTT
		},
	},
	"off": {
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpAckFrequency",
			Type:  "REG_DWORD",
			Value: "2", // 正常确认频率
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpNoDelay",
			Type:  "REG_DWORD",
			Value: "0", // 启用Nagle算法
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TcpDelAckTicks",
			Type:  "REG_DWORD",
			Value: "2", // 正常延迟
		},
		{
			Path:  `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
			Name:  "TCPInitialRTT",
			Type:  "REG_DWORD",
			Value: "3", // 正常RTT
		},
	},
}
