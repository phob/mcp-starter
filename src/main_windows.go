//go:build windows
// +build windows

package main

import (
	"os/exec"
	"syscall"
)

func executeServer(server McpServer) {
	cmd := exec.Command(server.Command, server.Args...)
	
	// Set up environment variables
	cmd.Env = os.Environ()
	for k, v := range server.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	// Windows-specific: Hide the command prompt window
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Connect stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic("Failed to run command: " + err.Error())
	}
} 