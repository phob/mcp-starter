package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type McpServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type Config struct {
	McpServers map[string]McpServer `json:"mcpServers"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run mcp-starter <config-file>")
		return
	}

	configPath := os.Args[1]
	config := loadConfig(configPath)
	
	// Get the single server entry
	if len(config.McpServers) != 1 {
		panic("Config must contain exactly one server definition")
	}
	
	var server McpServer
	var serverName string
	for name, s := range config.McpServers {
		serverName = name
		server = s
		break // Only need the first/only entry
	}

	fmt.Printf("Starting server: %s\n", serverName)
	executeServer(server)
}

func loadConfig(configPath string) Config {
	file, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error opening config file '%s': %v", configPath, err))
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Error decoding config '%s': %v", configPath, err))
	}

	return config
}

func executeServer(server McpServer) {
	cmd := exec.Command(server.Command, server.Args...)
	
	// Set up environment variables
	cmd.Env = os.Environ()
	for k, v := range server.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Windows-specific: Hide the command prompt window
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Connect stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("Failed to run command: %v", err))
	}
} 