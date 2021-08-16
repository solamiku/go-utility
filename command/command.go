package command

import (
	"os/exec"
	"runtime"

	"github.com/axgle/mahonia"
)

var Default *Command = &Command{}

func NewCommand() *Command {
	return &Command{
		envs: make(map[string]string, 2),
	}
}

type Command struct {
	dirname string
	envs    map[string]string
	args    []string
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

func (pcmd *Command) SetArgs(ps []string) {
	pcmd.args = ps
}

func (pcmd *Command) SetDir(dir string) {
	pcmd.dirname = dir
}

func (pcmd *Command) SetEnv(key, value string) {
	pcmd.envs[key] = value
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
	if len(pcmd.args) > 0 {
		args = append(args, pcmd.args...)
	}
	c := exec.Command(cmd, args...) //(cmd, args...)
	c.Dir = pcmd.dirname
	if len(pcmd.envs) > 0 {
		envarr := make([]string, 0, len(pcmd.envs))
		for k, val := range pcmd.envs {
			envarr = append(envarr, k+"="+val)
		}
		c.Env = envarr
	}
	ret, err := c.CombinedOutput()
	if err != nil {
		return string(ret), err
	}
	// var out bytes.Buffer
	// var errOut bytes.Buffer
	// c.Stdout = &out
	// c.Stderr = &errOut
	// err := c.Start()
	// if err != nil {
	// 	return errOut.String(), err
	// }
	// err = c.Wait()
	// if err != nil {
	// 	return errOut.String(), err
	// }
	// ret := out.String()
	return string(ret), err
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
