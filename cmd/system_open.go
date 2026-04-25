package cmd

import (
	"fmt"
	"os/exec"
)

func openSystemPath(target string) error {
	openCmd, args := getSystemOpenCommand()
	args = append(args, target)

	cmd := exec.Command(openCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法打开路径: %v (%s)", err, string(output))
	}
	return nil
}
