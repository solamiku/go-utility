package simplebar

import (
	"testing"
	"time"
)

func Test_simplebar(t *testing.T) {
	mgr := NewManager(500)
	a := mgr.NewSimpleBar(100, 0, "测试1")
	b := mgr.NewSimpleBar(100, 0, "测试2")
	go func() {
		for {
			select {
			case <-time.NewTicker(time.Second / 1000).C:
				a.AddNum(1)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.NewTicker(time.Second / 2).C:
				b.AddNum(1)
			}
		}
	}()
	mgr.Run()
	mgr.Wait()
}
