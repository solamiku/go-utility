package command

import (
	"bytes"
	"os/exec"
	"runtime"

	"github.com/axgle/mahonia"
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
	var errOut bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errOut
	err := c.Start()
	if err != nil {
		return encodeCommand(errOut.String()), err
	}
	err = c.Wait()
	if err != nil {
		return encodeCommand(errOut.String()), err
	}
	return encodeCommand(out.String()), err
}

func encodeCommand(src string) string {
	return ConvertToByte(src, "gbk", "utf8")
}

func ConvertToByte(src string, srcCode string, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}
