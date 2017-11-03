package utility

import (
	"sort"
	"strings"
	"time"
)

const (
	TIME_FMT_STR = "2006-01-02 15:04:05"
)

//return the golang time format string
//also you can use self custom format
//yyyy - year
//MM - month
//dd - day
//hh - hour
//mm - minutes
//ss - seconds
func TimeFmtStr(fmt ...string) string {
	if len(fmt) == 0 {
		return TIME_FMT_STR
	}
	replace := strings.NewReplacer(
		"yyyy", "2006",
		"MM", "01",
		"dd", "02",
		"hh", "15",
		"mm", "04",
		"ss", "05",
	)
	return replace.Replace(fmt[0])
}

// return today's specify time with hour,mintue,second
func TimeTodayCreate(h, m, s int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location())
}

//conver string to time
func TimeFromString(str string, fmt ...string) time.Time {
	fmtstr := TIME_FMT_STR
	if len(fmt) > 0 {
		fmtstr = fmt[0]
	}
	t, err := time.ParseInLocation(fmtstr, str, time.Local)
	if err != nil {
		return time.Unix(0, 0)
	}
	return t
}

//return the daily refresh hour time is expired or not;
//return next expired time reamin seconds;
func CheckDailyTimeExpire(now, lastTime time.Time, refreshHour []int) (bool, int) {
	t := 0
	t1 := 0
	if len(refreshHour) > 0 {
		t = refreshHour[0]
		t1 = refreshHour[0]
	}
	sort.Ints(refreshHour)
	for k := 1; k < len(refreshHour); k++ {
		if now.Hour() >= refreshHour[k] {
			t = refreshHour[k]
		}
		if now.Hour() < refreshHour[k] {
			t1 = refreshHour[k]
			break
		}
		if k == len(refreshHour)-1 {
			t1 = refreshHour[0]
			break
		}
	}
	timeFlash := TimeTodayCreate(t, 0, 0)
	if (lastTime.Before(timeFlash) && !now.Before(timeFlash)) ||
		lastTime.Before(timeFlash.AddDate(0, 0, -1)) {
		return true, 0
	}
	timeNext := TimeTodayCreate(t1, 0, 0)
	if now.Hour() >= t1 {
		timeNext = timeNext.AddDate(0, 0, 1)
	}
	duration := timeNext.Sub(now)
	return false, int(duration.Seconds() + 0.5)
}
