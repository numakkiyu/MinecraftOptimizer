package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// MenuHandler 菜单处理接口
type MenuHandler interface {
	Handle() error
}

// MainMenuHandler 主菜单处理
type MainMenuHandler struct{}

func (h *MainMenuHandler) Handle() error {
	for {
		showMainMenu()
		choice := getUserInput()

		var handler MenuHandler
		switch choice {
		case "1":
			handler = &BasicOptimizationHandler{}
		case "2":
			handler = &AdvancedOptimizationHandler{}
		case "3":
			handler = &ApplySettingsHandler{}
		case "4":
			handler = &BackupHandler{}
		case "0":
			return handleExit()
		default:
			fmt.Println("无效选择，请重试")
			continue
		}

		if err := handler.Handle(); err != nil {
			fmt.Printf("错误: %v\n", err)
		}

		fmt.Println("\n按Enter继续...")
		fmt.Scanln()
	}
}

// BasicOptimizationHandler 基础优化处理
type BasicOptimizationHandler struct{}

func (h *BasicOptimizationHandler) Handle() error {
	for {
		showBasicMenu()
		choice := getUserInput()

		switch choice {
		case "1": // 开启优化
			fmt.Println("正在优化网络设置...")
			mc, err := getMinecraftProcess()
			if err != nil {
				fmt.Printf("警告: 获取Minecraft进程失败: %v\n", err)
			} else {
				optimizer := NewNetworkOptimizer(mc)
				if err := optimizer.ApplyNetworkOptimizations(); err != nil {
					fmt.Printf("警告: 网络优化应用失败: %v\n", err)
				}
			}
			return nil

		case "2": // 关闭优化
			fmt.Println("正在恢复默认设置...")
			if err := disableOptimization(); err != nil {
				return fmt.Errorf("恢复设置失败: %v", err)
			}
			fmt.Println("已恢复默认设置！")
			return nil

		case "3": // 切换模式
			showModeMenu()
			modeChoice := getUserInput()
			switch modeChoice {
			case "1": // 命中优化
				return applyProfile(pvpProfiles["hit_reg"])
			case "2": // 击退优化
				return applyProfile(pvpProfiles["knockback"])
			case "3": // 连接稳定性
				return applyProfile(tuningProfiles["stability"])
			case "4": // 叠刀优化
				return handleStackedHit(true)
			}

		case "0":
			return nil
		default:
			fmt.Println("无效选择")
		}
	}
}

func applyProfile(profile NetworkProfile) error {
	fmt.Printf("正在应用 %s 配置...\n", profile.Name)
	fmt.Printf("说明: %s\n", profile.Description)

	// 应用注册表参数
	for _, cmd := range profile.RegParams {
		if err := executeRegCommand(cmd); err != nil {
			return fmt.Errorf("应用注册表参数失败: %v", err)
		}
	}

	// 应用网络参数
	for _, cmd := range profile.NetParams {
		if err := executeNetshCommand(cmd); err != nil {
			return fmt.Errorf("应用网络参数失败: %v", err)
		}
	}

	fmt.Printf("配置应用成功！\n")
	return nil
}

// AdvancedOptimizationHandler 高级优化处理
type AdvancedOptimizationHandler struct{}

func (h *AdvancedOptimizationHandler) Handle() error {
	for {
		showAdvancedMenu()
		choice := getUserInput()

		switch choice {
		case "1": // 网络调优
			return h.handleTuning()
		case "2": // 开启叠刀
			return handleStackedHit(true)
		case "3": // 关闭叠刀
			return handleStackedHit(false)
		case "0":
			return nil
		default:
			fmt.Println("无效选择")
		}
	}
}

func (h *AdvancedOptimizationHandler) handleTuning() error {
	showTuningMenu()
	choice := getUserInput()

	var profile string
	switch choice {
	case "1":
		profile = "disabled"
	case "2":
		profile = "default"
	case "3":
		profile = "restricted"
	case "4":
		profile = "delayed"
	case "0":
		return nil
	default:
		return fmt.Errorf("无效选择")
	}

	selectedProfile := tuningProfiles[profile]
	for _, cmd := range selectedProfile.NetParams {
		args := []string{
			"int",
			cmd.Category,
			"set",
			"global",
			fmt.Sprintf("%s=%s", cmd.Name, cmd.Value),
		}
		command := exec.Command("netsh", args...)
		if err := command.Run(); err != nil {
			return fmt.Errorf("设置网络参数失败: %v", err)
		}
	}
	return nil
}

// ApplySettingsHandler 应用设置处理
type ApplySettingsHandler struct{}

func (h *ApplySettingsHandler) Handle() error {
	// ps: 不会写
	fmt.Println("功能开发中...")
	return nil
}

// BackupHandler 备份处理
type BackupHandler struct{}

func (h *BackupHandler) Handle() error {
	for {
		showBackupMenu()
		choice := getUserInput()

		switch choice {
		case "1": // 创建备份
			if err := backupCurrentSettings(); err != nil {
				return fmt.Errorf("创建备份失败: %v", err)
			}
			fmt.Println("备份创建成功")
			return nil
		case "2": // 恢复备份
			if err := restoreFromBackup(); err != nil {
				return fmt.Errorf("恢复备份失败: %v", err)
			}
			fmt.Println("备份恢复成功")
			return nil
		case "3": // 查看备份
			if err := listBackups(); err != nil {
				return fmt.Errorf("查看备份失败: %v", err)
			}
			return nil
		case "0": // 返回主菜单
			return nil
		default:
			fmt.Println("无效选择")
		}
	}
}

// handleExit 程序退出处理
func handleExit() error {
	if currentStatus.IsEnabled {
		fmt.Println("正在恢复默认设置...")
		if err := disableOptimization(); err != nil {
			return fmt.Errorf("恢复设置失败: %v", err)
		}
	}
	return nil
}

// 叠刀设置处理
func handleStackedHit(enable bool) error {
	var commands []RegCommand
	if enable {
		commands = stackedHitSettings["on"]
		fmt.Println("正在开启叠刀优化...")
	} else {
		commands = stackedHitSettings["off"]
		fmt.Println("正在关闭叠刀优化...")
	}

	for _, cmd := range commands {
		if err := executeRegCommand(cmd); err != nil {
			return fmt.Errorf("设置叠刀参数失败: %v", err)
		}
	}

	// 额外优化
	if enable {
		netshCommands := []string{
			"int tcp set global autotuninglevel=restricted",
			"int tcp set global ecncapability=disabled",
			"int tcp set global timestamps=disabled",
		}

		for _, cmdStr := range netshCommands {
			cmd := exec.Command("netsh", strings.Split(cmdStr, " ")...)
			if err := cmd.Run(); err != nil {
				fmt.Printf("警告: 执行命令 '%s' 失败: %v\n", cmdStr, err)
			}
		}
	}

	fmt.Printf("叠刀优化已%s\n", map[bool]string{true: "开启", false: "关闭"}[enable])
	return nil
}
