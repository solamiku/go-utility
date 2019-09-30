package command

import "testing"

func Test_command(t *testing.T) {
	ret, err := Default.Run("date & dir", "/t")
	if err != nil {
		t.Fatalf("run err:%v", err)
	} else {
		t.Log(ret)
	}

	cmd := NewCommand()
	cmd.SetDecode("gbk", "utf8")
	ret, err = cmd.Run("date & dir", "/t")
	if err != nil {
		t.Fatalf("run err:%v", err)
	} else {
		t.Log(ret)
	}
}
