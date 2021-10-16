package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, []int{3, 1, 2}, res, "Call order")
}
