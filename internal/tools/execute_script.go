package tools

import (
	"fmt"
	"os/exec"
)

// ExecuteScript calls a script with positional parameters
func ExecuteScript(path string, args string) string {
	cmd, err := exec.Command("/bin/sh", path, args).CombinedOutput()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string([]byte(cmd))
	return output
}
