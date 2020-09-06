package task

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestManager(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)
	res := []int64{0, 0, 0, 0, 0}

	m := NewManager()
	m.
		Register("task1", Low(func() {
			atomic.AddInt64(&res[1], 1)
			m.Schedule("finish")
		})).
		Register("task2", Medium(func() {
			atomic.AddInt64(&res[2], 1)
			m.Schedule("task1")
		})).
		Register("task3", High(func() {
			atomic.AddInt64(&res[3], 1)
			m.Schedule("task1").Schedule("task2")
		})).
		Register("finish", Idle(wg.Done)).
		Start()

	m.
		ScheduleAction(func() {
			atomic.AddInt64(&res[0], 1)
			m.Schedule("task3")
		}).
		ScheduleTask(Realtime(func() {
			atomic.AddInt64(&res[4], 1)
			m.Schedule("task3")
		})).
		Schedule("finish")

	wg.Wait()
	needRes := []int64{1, 4, 2, 2, 1}
	if !reflect.DeepEqual(needRes, res) {
		t.Errorf("Incorrect count of Task calls. Want %v, got %v", needRes, res)
	}
}
