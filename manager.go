package task

import (
	"sync"
)

type Manager struct {
	mu        sync.RWMutex
	queue     []queue
	taskStore map[string]Task
	play      chan struct{}
}

func NewManager() *Manager {
	queues := make([]queue, total)
	for i := maxPriority; i >= 0; i-- {
		queues[i] = queue{}
	}
	return (&Manager{
		queue:     queues,
		taskStore: map[string]Task{},
		play:      make(chan struct{}),
	}).init()
}

func (m *Manager) Register(name string, task Task) *Manager {
	m.mu.Lock()
	m.taskStore[name] = task
	m.mu.Unlock()
	return m
}

func (m *Manager) scheduleTask(task Task) {
	m.mu.Lock()
	m.queue[task.priority].Enqueue(task.Action)
	m.mu.Unlock()
	go m.next()
}

func (m *Manager) getNextAction() Action {
	var task Action
	m.mu.Lock()
	for i := maxPriority; i >= 0 && task == nil; i-- {
		task = m.queue[i].Dequeue()
	}
	m.mu.Unlock()
	return task
}

func (m *Manager) Schedule(name string) *Manager {
	m.mu.RLock()
	task, ok := m.taskStore[name]
	m.mu.RUnlock()
	if ok {
		m.scheduleTask(task)
	}
	return m
}

func (m *Manager) ScheduleTask(task Task) *Manager {
	m.scheduleTask(task)
	return m
}

func (m *Manager) ScheduleAction(action Action) *Manager {
	m.scheduleTask(Task{
		Action:   action,
		priority: idle,
	})
	return m
}

func (m *Manager) next() {
	m.play <- struct{}{}
}

func (m *Manager) processTasks() {
	for {
		if task := m.getNextAction(); task != nil {
			task()
		} else {
			return
		}
	}
}

func (m *Manager) loop() {
	for {
		m.processTasks()
		<-m.play
	}
}

func (m *Manager) init() *Manager {
	go m.loop()
	return m
}
