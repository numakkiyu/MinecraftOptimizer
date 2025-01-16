package main

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

// 读取注册表值
func readRegistryValue(keyName string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("打开注册表键失败: %v", err)
	}
	defer key.Close()

	// 读取值
	value, _, err := key.GetStringValue(keyName)
	if err != nil {
		// 如果值不存在，返回默认值
		if err == registry.ErrNotExist {
			return "0", nil
		}
		return "", fmt.Errorf("读取注册表值失败: %v", err)
	}

	return value, nil
}

// 写入注册表值
func writeRegistryValue(keyName, valueType, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %v", err)
	}
	defer key.Close()

	// 写入值
	err = key.SetStringValue(keyName, value)
	if err != nil {
		return fmt.Errorf("写入注册表值失败: %v", err)
	}

	return nil
}

// 备份注册表键
func backupRegistryKey(keyPath string) error {
	// 使用reg export命令备份
	cmd := exec.Command("reg", "export", keyPath, backupPath, "/y")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("备份注册表失败: %v, 输出: %s", err, output)
	}
	return nil
}

// 恢复注册表键
func restoreRegistryKey(backupFile string) error {
	// 使用reg import命令恢复
	cmd := exec.Command("reg", "import", backupFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("恢复注册表失败: %v, 输出: %s", err, output)
	}
	return nil
}
