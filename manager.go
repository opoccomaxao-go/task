package task

import (
	"reflect"
	"sync"
)

type Manager struct {
	mu        sync.RWMutex
	queue     []chan Action
	cases     []reflect.SelectCase
	taskStore map[string]Task
	canStart  bool
}

func NewManager() *Manager {
	return NewManagerCap(10)
}

func NewManagerCap(capacity int) *Manager {
	queue := make([]chan Action, total)
	cases := make([]reflect.SelectCase, total)
	if capacity < 1 {
		capacity = 1
	}
	for i := total - 1; i >= 0; i-- {
		ch := make(chan Action, capacity)
		queue[i] = ch
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}
	return &Manager{
		queue:     queue,
		cases:     cases,
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
		m.queue[task.priority] <- task.Action
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
func (m *Manager) ScheduleAction(action Action) *Manager {
	m.scheduleTask(Task{
		Action:   action,
		priority: idle,
	})
	return m
}

func (m *Manager) loop() {
	for {
		_, task, _ := reflect.Select(m.cases)
		task.Interface().(Action)()
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
