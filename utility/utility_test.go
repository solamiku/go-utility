package utility

import (
	"testing"
	"time"
)

func Test_utility(t *testing.T) {
	now := time.Now()
	t.Log(now.Format(TimeFmtStr()))
	t.Log(now.Format(TimeFmtStr("yyyy hhh")))
	t.Log(TimeFromString("2017-11-02 23:00:00"))

	//test check daily time expire
	type checkDailysS struct {
		now     string
		last    string
		refresh []int
		f       bool
		remain  int
	}
	checkDailys := []checkDailysS{
		checkDailysS{"2017-11-03 04:01:00", "2017-11-02 23:00:00", []int{4}, true, 0},
		checkDailysS{"2017-11-03 03:00:00", "2017-11-02 23:00:00", []int{4}, false, 3600},
		checkDailysS{"2017-11-03 03:00:00", "2017-11-03 02:00:00", []int{0, 2}, false, 75600},
		checkDailysS{"2017-11-03 01:00:00", "2017-11-03 01:00:00", []int{0, 2}, false, 3600},
		checkDailysS{"2017-11-03 03:00:00", "2017-11-03 01:00:00", []int{0, 2}, true, 0},
	}
	for _, v := range checkDailys {
		f, r := CheckDailyTimeExpire(TimeFromString(v.now), TimeFromString(v.last), v.refresh)
		if f != v.f || r != v.remain {
			t.Fatalf("%v checkdaily failed. return is %v %d", v, f, r)
		}
	}
	t.Log("test check daily done.")

	//test utility
	t.Log("test utility")
	t.Log(Str(now))
}
