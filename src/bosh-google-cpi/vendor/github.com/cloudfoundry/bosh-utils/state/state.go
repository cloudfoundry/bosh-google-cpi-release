package state

import (
	"errors"
	"sync/atomic"
)

type Statuser interface {
	String() string
	Int32() int32 // The returned int32 must be static.
}

type Status int32

const (
	Stopped Status = iota
	Running
	Exited
)

var statusStr = [...]string{
	"Stopped",
	"Running",
	"Exited",
}

func (s Status) String() string { return statusStr[s] }
func (s Status) Int32() int32   { return int32(s) }

type State struct {
	i int32
	m map[int32]Statuser
}

// Returns a new State for the provided Statusers.  The Status of the
// returned State is set to the first Statuser argument.
func NewState(sts ...Statuser) (*State, error) {
	if len(sts) < 2 {
		return nil, errors.New("state: must provide at least 2 Statuser arguments")
	}
	m := make(map[int32]Statuser, len(sts))
	for _, s := range sts {
		if _, ok := m[s.Int32()]; ok {
			return nil, errors.New("state: duplicate state")
		}
		m[s.Int32()] = s
	}
	return &State{i: sts[0].Int32(), m: m}, nil
}

// Is returns if the State has status st.
func (s *State) Is(st Statuser) bool {
	return atomic.LoadInt32(&s.i) == s.value(st)
}

// Set, sets the State to status st.
func (s *State) Set(st Statuser) {
	atomic.StoreInt32(&s.i, s.value(st))
}

// Status returns the current status.
func (s *State) Status() Statuser {
	return s.m[atomic.LoadInt32(&s.i)]
}

// Swap, sets the status to new and returns the previous status.
func (s *State) Swap(new Statuser) (old Statuser) {
	return s.m[atomic.SwapInt32(&s.i, s.value(new))]
}

// CompareAndSwap, executes the compare-and-swap operation for the status.
func (s *State) CompareAndSwap(old, new Statuser) (swapped bool) {
	return atomic.CompareAndSwapInt32(&s.i, s.value(old), s.value(new))
}

func (s *State) value(st Statuser) int32 {
	i := st.Int32()
	if _, ok := s.m[i]; !ok {
		panic("state: invalid Statuser")
	}
	return i
}
