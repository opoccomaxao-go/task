package task

type Task struct {
	Action
	priority
}

func Idle(action Action) Task {
	return Task{
		Action:   action,
		priority: idle,
	}
}

func Low(action Action) Task {
	return Task{
		Action:   action,
		priority: low,
	}
}

func Medium(action Action) Task {
	return Task{
		Action:   action,
		priority: medium,
	}
}

func High(action Action) Task {
	return Task{
		Action:   action,
		priority: high,
	}
}

func Realtime(action Action) Task {
	return Task{
		Action:   action,
		priority: realtime,
	}
}
