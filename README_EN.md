# Minecraft Network Optimizer

<div align="center">

![Version](https://img.shields.io/badge/version-0.2.0-blue.svg)
![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)
![Platform](https://img.shields.io/badge/platform-Windows-lightgrey.svg)
![Language](https://img.shields.io/badge/language-Go-00ADD8.svg)

A professional Minecraft network optimization tool, focusing on enhancing the PVP experience and network performance.

English | [Simplified Chinese](./README_EN.md)

</div>

## 📝 Table of Contents

- [Features](#-features)
- [Technical Implementation](#-technical-implementation)
- [Installation Instructions](#-installation-instructions)
- [User Guide](#-user-guide)
- [Compilation Guide](#-compilation-guide)
- [Contribution Guide](#-contribution-guide)
- [Open Source License](#-open-source-license)

## ✨ Features

### Basic Optimization
- Hit Optimization: Improve hit judgment and feedback.
- Knockback Optimization: Optimize knockback effects and synchronization.
- Combo Optimization: Optimize combo judgment.
- Network Optimization: Adaptive network tuning.

### Advanced Features
- QoS Policy Management
- TCP/UDP Parameter Optimization
- Network Buffer Adjustment
- System Parameter Optimization
- Configuration Backup and Restoration

## 🛠 Technical Implementation

### Core Technology Stack
- Language: Go 1.21+
- System: Windows
- Dependencies:
  - `golang.org/x/sys`
  - Windows API
  - PowerShell

### Main Modules
```
├── network/        # Core of network optimization
├── process/        # Process management
├── registry/       # Registry operations
├── profiles/       # Optimization configuration files
└── ui/            # User interface
```

### Optimization Principles
1. TCP Parameter Optimization
   - TcpAckFrequency
   - TcpNoDelay
   - TCPInitialRTT
   - DefaultTTL

2. QoS Policy
   - Application Priority
   - Network Packet Priority
   - Bandwidth Management

3. Network Buffer
   - Adaptive Tuning
   - Buffer Size Optimization
   - Latency Control

## 📦 Installation Instructions

### System Requirements
- Windows 10/11
- Administrator Privileges
- .NET Framework 4.5+

### Download and Install
1. Download the latest version from [Releases](https://github.com/numakkiyu/MinecraftOptimizer/releases).
2. Run the program as an administrator.
3. Follow the prompts to proceed.

## 🚀 Compilation Guide

### Environment Preparation
1. Install Go 1.21 or a higher version.
2. Install Git.
3. Install rcedit (for resource file editing).

### Compilation Steps
```bash
# Clone the repository
git clone https://github.com/numakkiyu/MinecraftOptimizer.git
cd MinecraftOptimizer

# Install dependencies
go mod tidy

# Compile
go build -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d'` -X main.GitCommit=`git rev-parse --short HEAD`" -o minecraft_optimizer.exe
```

## 📖 User Guide

### Basic Optimization
1. Start Minecraft.
2. Run the optimizer.
3. Select "Basic Optimization".
4. Select the optimization mode.

### Advanced Settings
- Network Tuning: Customize network parameters.
- Combo Optimization: Optimize combo effects.
- Configuration Management: Import and export configurations.


## 📞 Contact Information

- Blog: https://me.tianbeigm.cn
- GitHub: https://github.com/numakkiyu/MinecraftOptimizer
