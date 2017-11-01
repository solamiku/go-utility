package command

import (
	"bytes"
	"os/exec"
	"runtime"
)

func Run(cmd string, args ...string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return runCmd("cmd", append([]string{"/C", cmd}, args...)...)
	default:
		return runCmd(cmd, args...)
	}
}

func runCmd(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...) //(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Start()
	if err != nil {
		return "", err
	}
	err = c.Wait()
	if err != nil {
		return "", err
	}
	return out.String(), err
}