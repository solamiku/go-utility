package command

import (
	"bytes"
	"os/exec"
	"runtime"

	"github.com/axgle/mahonia"
)

var Default *Command = &Command{}

func NewCommand() *Command {
	return &Command{}
}

type Command struct {
	dirname string
	decode  []string
}

func (pcmd *Command) Run(cmd string, args ...string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return pcmd.runCmd("cmd", append([]string{"/C", cmd}, args...)...)
	default:
		return pcmd.runCmd(cmd, args...)
	}
}

func (pcmd *Command) SetDir(dir string) {
	pcmd.dirname = dir
}

func (pcmd *Command) SetDecode(src, dst string) {
	if len(src) == 0 {
		pcmd.decode = pcmd.decode[:0]
	}
	pcmd.decode = []string{src, dst}
}

//
func (pcmd *Command) runCmd(cmd string, args ...string) (rmsg string, rerr error) {
	defer func() {
		if len(pcmd.decode) > 0 {
			rmsg = encodeCommand(rmsg, pcmd.decode[0], pcmd.decode[1])
		}
	}()
	c := exec.Command(cmd, args...) //(cmd, args...)
	c.Dir = pcmd.dirname
	var out bytes.Buffer
	var errOut bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errOut
	err := c.Start()
	if err != nil {
		return errOut.String(), err
	}
	err = c.Wait()
	if err != nil {
		return errOut.String(), err
	}
	ret := out.String()
	return ret, err
}

func encodeCommand(src, s, t string) string {
	return ConvertToByte(src, s, t)
}

func ConvertToByte(src string, srcCode string, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}
