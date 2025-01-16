package main

import (
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
)

// checkAdminPrivileges 检查是否具有管理员权限
func checkAdminPrivileges() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

// 检查进程是否以管理员权限运行的备用方法
func checkAdminPrivilegesAlt() bool {
	cmd := exec.Command("net", "session")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return !strings.Contains(string(output), "拒绝访问")
}
