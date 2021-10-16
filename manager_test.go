package task

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)
	res := []int64{0, 0, 0, 0, 0}

	m := NewManager()
	m.
		Register("task1", Low(func() {
			res[1]++
			go m.Schedule("finish")
		})).
		Register("task2", Medium(func() {
			res[2]++
			go m.Schedule("task1")
		})).
		Register("task3", High(func() {
			res[3]++
			go m.Schedule("task1").Schedule("task2")
		})).
		Register("finish", Idle(func() {
			wg.Done()
		}))

	m.
		ScheduleAction(func() {
			res[0]++
			go m.Schedule("task3")
		}).
		ScheduleTask(Realtime(func() {
			res[4]++
			m.Schedule("task3")
		})).
		Schedule("finish")
	wg.Wait()

	assert.Equal(t, []int64{1, 4, 2, 2, 1}, res, "Call count")
}
