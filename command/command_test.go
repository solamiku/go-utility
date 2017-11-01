package command

import "testing"

func Test_command(t *testing.T) {
	ret, err := Run("date", "/t")
	if err != nil {
		t.Fatal("run err:%v", err)
	} else {
		t.Log(ret)
	}
}
