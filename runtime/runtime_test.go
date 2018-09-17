package runtime

import (
	"bytes"
	"os"
	"testing"
)

func Test_runtime(t *testing.T) {
	t.Logf("the number of heaps is : %d", LookupHeapObjs())
	t.Log(Errof("test"))
	t.Log(CallInfo(1))

	a := bytes.NewBuffer([]byte(""))
	WriteRoutineCallstack(0, a)
	t.Log("call stack\n", a.String())

	t.Log(IsPahtExisted("log"))
	t.Log(IsPahtExisted("log/test"))
	t.Log(IsPahtExisted("log/test1"))
	t.Log(MakeDir("log/test", os.ModePerm))
	t.Log(IsPahtExisted("log/test"))
}
