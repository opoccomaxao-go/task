package task

import (
	"sync"
	"sync/atomic"
)

type Manager struct {
	mu        sync.RWMutex
	run       sync.Mutex
	queue     []queue
	taskStore map[string]Task
	active    int32
}

func NewManager() *Manager {
	queues := make([]queue, total)
	for i := maxPriority; i >= 0; i-- {
		queues[i] = queue{}
	}
	return &Manager{
		queue:     queues,
		taskStore: map[string]Task{},
	}
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
	defer atomic.AddInt32(&m.active, -1)
	if atomic.AddInt32(&m.active, 1) > 1 {
		return
	}
	for {
		task := m.getNextAction()
		if task == nil {
			return
		}
		task()
	}
}
