package tools

import (
	"log"
	"os/exec"
)

// ExecuteScript calls a script with positional parameters
func ExecuteScript(args ...string) string {
	cmd, err := exec.Command("/bin/sh", args...).CombinedOutput()
	if err != nil {
		log.Fatal()
	}
	output := string([]byte(cmd))
	return output
}
