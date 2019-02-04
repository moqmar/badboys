package main

import (
	"os/exec"
	"strings"
)

// uses net=host
func runInDocker(volumes []string, environment []string, command string, args ...string) *exec.Cmd {
	// TODO:
	return nil
}

func resolveIP(hostname string) string {
	if strings.HasPrefix(hostname, "docker=") {
		// TODO:
		return ""
	}
	return hostname
}
