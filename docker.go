package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// uses net=host
func runInDocker(volumes []string, environment []string, command string, args ...string) *exec.Cmd {
	if cfg.Get("localTools").Bool() {
		cmd := exec.Command(command, args...)
		// Replace volume paths
		for _, v := range volumes {
			vp := strings.SplitN(v, ":", 2)
			command = strings.Replace(command, vp[1], vp[0], -1)
			for i := range args {
				args[i] = strings.Replace(args[i], vp[1], vp[0], -1)
			}
			for i := range environment {
				environment[i] = strings.Replace(environment[i], vp[1], vp[0], -1)
			}
		}
		cmd.Env = append(os.Environ(), environment...)
		return cmd
	}

	docker := []string{"run", "--net=host", "--rm"}
	for _, v := range volumes {
		docker = append(docker, "-v", v)
	}
	for _, e := range environment {
		docker = append(docker, "-e", e)
	}
	docker = append(append(docker, command), args...)
	return exec.Command("docker", docker...)

}

func resolveIP(hostname string) string {
	if strings.HasPrefix(hostname, "docker=") {
		out, err := exec.Command("docker", "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", strings.SplitN(hostname, "=", 2)[1]).Output()
		if err != nil {
			fmt.Printf("Couldn't resolve docker hostname %s: %s", strings.SplitN(hostname, "=", 2)[1], err)
			return hostname
		}
		return strings.SplitN(string(out), "\n", 2)[0]
	}
	return hostname
}
