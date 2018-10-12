package simplebar

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type SimpleBar struct {
	id      int
	max     int
	cur     int
	preInfo string

	barMutex   sync.Mutex
	barLength  int
	maxTimeout int
	startTime  time.Time
	finishStr  string
}

func (sb *SimpleBar) CheckFinish(tNow time.Time) bool {
	if sb.cur >= sb.max {
		return true
	}
	if sb.maxTimeout > 0 {
		if int(tNow.Sub(sb.startTime).Seconds()*1000) > sb.maxTimeout {
			return true
		}
	}
	return false
}

func (sb *SimpleBar) AddNum(num int) {
	sb.cur += num
}

func (sb *SimpleBar) SetNum(num int) {
	sb.cur = num
}

//

type SimpleBarManager struct {
	id           int
	barsMutex    sync.Mutex
	bars         map[int]*SimpleBar
	savePos      bool
	interval     int
	finishNotify chan int
}

func (sm *SimpleBarManager) Run() {
	go func() {
		defer func() {
			sm.finishNotify <- 1
		}()
		msecTimer := time.NewTicker(time.Duration(sm.interval) * time.Millisecond)
		for {
			select {
			case <-msecTimer.C:
				//output
				sm.PrintAll()
				//detect
				if sm.CheckAllFinish() {
					return
				}
			}
		}
	}()
}

func (sm *SimpleBarManager) Wait() {
	if sm.finishNotify == nil {
		return
	}
	<-sm.finishNotify
}

func (sm *SimpleBarManager) NewSimpleBar(max, timeout int, preInfo string) *SimpleBar {
	sm.barsMutex.Lock()
	defer sm.barsMutex.Unlock()
	sm.id++
	bar := &SimpleBar{
		id:         0,
		max:        max,
		cur:        0,
		preInfo:    preInfo,
		barLength:  35,
		startTime:  time.Now(),
		maxTimeout: timeout,
	}
	bar.id = sm.id
	sm.bars[sm.id] = bar
	return bar
}

func (sm *SimpleBarManager) PrintAll() {
	sm.barsMutex.Lock()
	defer sm.barsMutex.Unlock()
	now := time.Now()
	fmt.Print("\r")
	if !sm.savePos {
		sm.savePos = true
		fmt.Print("\033[s")
	}
	fmt.Print("\033[u")
	for i := 1; i <= sm.id; i++ {
		bar, ok := sm.bars[i]
		if !ok {
			continue
		}
		barString := bar.finishStr
		if len(barString) == 0 {
			curl := int(float64(bar.cur) / float64(bar.max) * float64(bar.barLength))
			if curl > bar.barLength {
				curl = bar.barLength
			}
			remainl := bar.barLength - curl
			timeSec := int(now.Sub(bar.startTime).Seconds())
			timeFmt := fmt.Sprintf(
				"%dday %02d:%02d:%02d",
				timeSec/3600/24,
				timeSec/3600%24,
				timeSec/60%60,
				timeSec%60,
			)
			speed := float64(bar.max)
			if timeSec != 0 {
				speed = float64(bar.cur) / float64(timeSec)
			}
			barString = fmt.Sprintf(
				"[%s](%v) {%s%s} [%d/%d %d%% speed:%.3f/s]\n",
				timeFmt,
				bar.preInfo,
				strings.Repeat("*", curl),
				strings.Repeat("-", remainl),
				bar.cur,
				bar.max,
				bar.cur*100/bar.max,
				speed,
			)
		}
		fmt.Print(barString)
		if bar.CheckFinish(now) {
			bar.finishStr = barString
		}
	}
}

func (sm *SimpleBarManager) CheckAllFinish() bool {
	sm.barsMutex.Lock()
	defer sm.barsMutex.Unlock()
	now := time.Now()
	for _, bar := range sm.bars {
		if !bar.CheckFinish(now) {
			return false
		}
	}
	return true
}

//

func NewManager(interval int) *SimpleBarManager {
	barManagers := &SimpleBarManager{}
	barManagers.bars = make(map[int]*SimpleBar)
	barManagers.finishNotify = make(chan int, 1)
	barManagers.interval = interval
	fmt.Println("!!!Notice:simple processbar running!!!")
	fmt.Println("!!!Any other output will be clear per interval!!!")
	return barManagers

}
