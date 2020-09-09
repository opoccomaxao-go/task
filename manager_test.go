package task

import (
	"fmt"
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
			fmt.Println("task1")
			atomic.AddInt64(&res[1], 1)
			go m.Schedule("finish")
		})).
		Register("task2", Medium(func() {
			fmt.Println("task2")
			atomic.AddInt64(&res[2], 1)
			go m.Schedule("task1")
		})).
		Register("task3", High(func() {
			fmt.Println("task3")
			atomic.AddInt64(&res[3], 1)
			go m.Schedule("task1").Schedule("task2")
		})).
		Register("finish", Idle(func() {
			fmt.Println("finish")
			wg.Done()
		}))

	m.
		ScheduleAction(func() {
			fmt.Println("task0")
			atomic.AddInt64(&res[0], 1)
			go m.Schedule("task3")
		}).
		ScheduleTask(Realtime(func() {
			fmt.Println("task4")
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
