package task

type queue struct {
	head *node
	tail *node
}

type node struct {
	Action Action
	next   *node
}

func (q *queue) Enqueue(action Action) {
	n := &node{Action: action}
	if q.tail == nil {
		q.head = n
	} else {
		q.tail.next = n
	}
	q.tail = n
}

func (q *queue) Dequeue() Action {
	if q.head == nil {
		return nil
	}
	res := q.head.Action
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return res
}
