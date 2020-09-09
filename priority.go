package task

type priority int

const (
	idle priority = iota
	low
	medium
	high
	realtime
	total
)

const maxPriority = total - 1
