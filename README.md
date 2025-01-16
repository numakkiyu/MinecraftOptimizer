# Minecraft Network Optimizer

<div align="center">

![Version](https://img.shields.io/badge/version-0.2.0-blue.svg)
![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)
![Platform](https://img.shields.io/badge/platform-Windows-lightgrey.svg)
![Language](https://img.shields.io/badge/language-Go-00ADD8.svg)

一个专业的 Minecraft 网络优化工具，专注于提升 PVP 体验和网络性能。

[English](./README_EN.md) | 简体中文

</div>

## 📝 目录

- [功能特性](#-功能特性)
- [技术实现](#-技术实现)
- [安装说明](#-安装说明)
- [使用指南](#-使用指南)
- [编译指南](#-编译指南)
- [贡献指南](#-贡献指南)
- [开源协议](#-开源协议)

## ✨ 功能特性

### 基础优化
- 命中优化：提升击打判定和命中反馈
- 击退优化：优化击退效果和同步性
- 叠刀优化：优化连击判定
- 网络优化：自适应网络调优

### 高级功能
- QoS 策略管理
- TCP/UDP 参数优化
- 网络缓冲区调整
- 系统参数优化
- 配置备份与恢复

## 🛠 技术实现

### 核心技术栈
- 语言：Go 1.21+
- 系统：Windows
- 依赖：
  - `golang.org/x/sys`
  - Windows API
  - PowerShell

### 主要模块
```
├── network/        # 网络优化核心
├── process/        # 进程管理
├── registry/       # 注册表操作
├── profiles/       # 优化配置文件
└── ui/            # 用户界面
```

### 优化原理
1. TCP 参数优化
   - TcpAckFrequency
   - TcpNoDelay
   - TCPInitialRTT
   - DefaultTTL

2. QoS 策略
   - 应用程序优先级
   - 网络包优先级
   - 带宽管理

3. 网络缓冲区
   - 自适应调优
   - 缓冲区大小优化
   - 延迟控制

## 📦 安装说明

### 系统要求
- Windows 10/11
- 管理员权限
- .NET Framework 4.5+

### 下载安装
1. 从 [Releases](https://github.com/numakkiyu/MinecraftOptimizer/releases) 下载最新版本
2. 以管理员身份运行程序
3. 按照提示进行操作

## 🚀 编译指南

### 环境准备
1. 安装 Go 1.21 或更高版本
2. 安装 Git
3. 安装 rcedit (用于资源文件编辑)

### 编译步骤
```bash
# 克隆仓库
git clone https://github.com/numakkiyu/MinecraftOptimizer.git
cd MinecraftOptimizer

# 安装依赖
go mod tidy

# 编译
go build -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d'` -X main.GitCommit=`git rev-parse --short HEAD`" -o minecraft_optimizer.exe
```

## 📖 使用指南

### 基础优化
1. 启动 Minecraft
2. 运行优化器
3. 选择"基础优化"
4. 选择优化模式

### 高级设置
- 网络调优：自定义网络参数
- 叠刀优化：优化连击效果
- 配置管理：导入导出配置


## 📞 联系方式

- 博客：https://me.tianbeigm.cn
- GitHub：https://github.com/numakkiyu/MinecraftOptimizer

