package task

import (
	"sync"
)

type Manager struct {
	mu        sync.RWMutex
	queue     chan Task
	taskStore map[string]Task
	canStart  bool
}

func NewManager() *Manager {
	return &Manager{
		queue:     make(chan Task, 100),
		taskStore: map[string]Task{},
		canStart:  true,
	}
}

func (m *Manager) Register(name string, task Task) *Manager {
	m.mu.Lock()
	m.taskStore[name] = task
	m.mu.Unlock()
	return m
}

func (m *Manager) scheduleTask(task Task) {
	go func() {
		m.queue <- task
	}()
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

func (m *Manager) loop() {
	for task := range m.queue {
		task()
	}
}

func (m *Manager) Start() {
	m.mu.Lock()
	if m.canStart {
		m.canStart = false
		go m.loop()
	}
	m.mu.Unlock()
}
