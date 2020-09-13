package task

import (
	"fmt"
	"reflect"
	"sync"
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
			res[1]++
			go m.Schedule("finish")
		})).
		Register("task2", Medium(func() {
			fmt.Println("task2")
			res[2]++
			go m.Schedule("task1")
		})).
		Register("task3", High(func() {
			fmt.Println("task3")
			res[3]++
			go m.Schedule("task1").Schedule("task2")
		})).
		Register("finish", Idle(func() {
			fmt.Println("finish")
			wg.Done()
		}))

	m.
		ScheduleAction(func() {
			fmt.Println("task0")
			res[0]++
			go m.Schedule("task3")
		}).
		ScheduleTask(Realtime(func() {
			fmt.Println("task4")
			res[4]++
			m.Schedule("task3")
		})).
		Schedule("finish")
	wg.Wait()
	needRes := []int64{1, 4, 2, 2, 1}
	if !reflect.DeepEqual(needRes, res) {
		t.Errorf("Incorrect count of Task calls. Want %v, got %v", needRes, res)
	}
}
