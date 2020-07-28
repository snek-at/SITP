package tools

func ExecuteScript(path string, pos_args string) string {
	cmd, err := exec.Command("/bin/sh", path, pos_args).CombinedOutput()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string([]byte(cmd))
	return output
}