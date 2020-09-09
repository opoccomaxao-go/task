package task

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	var res []int
	q := new(queue)
	q.Enqueue(func() {
		res = append(res, 3)
	})
	q.Enqueue(func() {
		res = append(res, 1)
	})
	q.Enqueue(func() {
		res = append(res, 2)
	})
	for {
		act := q.Dequeue()
		if act == nil {
			break
		} else {
			act()
		}
	}
	needRes := []int{3, 1, 2}
	if !reflect.DeepEqual(needRes, res) {
		t.Errorf("Incorrect call order. Want %v, got %v", needRes, res)
	}
}
