//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func executeServer(server McpServer) {
	cmd := exec.Command(server.Command, server.Args...)
	
	// Set up environment variables
	cmd.Env = os.Environ()
	for k, v := range server.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Connect stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("Failed to run command: %v", err))
	}
} 