package main

// RegCommand 注册表命令
type RegCommand struct {
	Path  string // 注册表路径
	Name  string // 键名
	Type  string // 值类型
	Value string // 值
}

// NetshCommand 网络命令
type NetshCommand struct {
	Category string // 类别
	Setting  string // 设置
	Name     string // 名称
	Value    string // 值
}

// NetworkProfile 网络配置文件
type NetworkProfile struct {
	Name        string         // 配置名称
	Description string         // 配置描述
	RegParams   []RegCommand   // 注册表参数
	NetParams   []NetshCommand // 网络参数
	BufferSize  int            // 缓冲区大小
	Priority    string         // 优先级
}

// MinecraftProcess MC进程信息
type MinecraftProcess struct {
	PID      string
	Name     string
	Title    string
	Memory   int64
	NetUsage int64
	Path     string // 进程完整路径
}
