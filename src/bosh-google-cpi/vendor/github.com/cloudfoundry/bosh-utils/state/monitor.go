package state

import (
	"runtime"
	"sync"
)

type Monitor struct {
	state *State
	c     *sync.Cond
	mu    sync.RWMutex
}

// New returns a new Monitor.
func New() (*Monitor, error) {
	state, err := NewState(Stopped, Running, Exited)
	if err != nil {
		return nil, err
	}
	m := &Monitor{
		state: state,
		c:     &sync.Cond{L: new(sync.Mutex)},
	}
	return m, nil
}

// NewProcess returns a new monitored Process.
func (m *Monitor) NewProcess() (p *Process) {
	m.mu.RLock()
	p = &Process{state: m.state, c: m.c}
	m.mu.RUnlock()
	return
}

func (m *Monitor) Status() (st Statuser) {
	m.mu.RLock()
	st = m.state.Status()
	m.mu.RUnlock()
	return
}

// Stop, stops exection of all monitored Processes.
func (m *Monitor) Stop() {
	m.mu.RLock()
	m.broadcast(Stopped)
	m.mu.RUnlock()
}

// Start, starts exection of all monitored Processes.
func (m *Monitor) Start() {
	m.mu.RLock()
	m.broadcast(Running)
	m.mu.RUnlock()
}

// Exit, causes all monitored Processes to exit and the status set to Stopped.
func (m *Monitor) Exit() {
	m.mu.Lock()
	m.broadcast(Exited)
	m.state, _ = NewState(Stopped, Running, Exited)
	m.c = &sync.Cond{L: new(sync.Mutex)}
	m.mu.Unlock()
}

func (m *Monitor) broadcast(s Statuser) {
	m.state.Set(s)
	m.c.Broadcast()
}

type Process struct {
	state *State
	c     *sync.Cond
}

// Wait is used in goroutines to respond to the Monitor State changes.
// Wait controls exectution of the calling goroutine, allowing the Monitor
// to stop, start and exit execution.  When Exit is called by the Monitor
// Wait calls runtime.Goexit halting execution of the calling goroutine.
//
// Example:
//
//    m := New()
//    p := m.NewProcess()
//    go func() {
//        defer cleanup() // m.Exit() calls runtime.Goexit()
//        for {
//            p.Wait()
//            ... do stuff ...
//        }
//    }()
//
func (p *Process) Wait() {
	if p.state.Is(Running) {
		return
	}
	if p.state.Is(Exited) {
		runtime.Goexit()
	}
	if p.state.Is(Stopped) {
		p.c.L.Lock()
		for p.state.Is(Stopped) {
			p.c.Wait()
		}
		p.c.L.Unlock()
		if !p.state.Is(Running) {
			p.Wait()
		}
	}
}
