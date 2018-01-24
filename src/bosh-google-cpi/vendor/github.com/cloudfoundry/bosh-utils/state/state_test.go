package state

import "testing"

type TestStatus int32

const (
	TestStopped TestStatus = iota
	TestRunning
	TestExited
)

var testStatusStr = [...]string{
	"TestStopped",
	"TestRunning",
	"TestExited",
}

func (s TestStatus) String() string { return testStatusStr[s] }
func (s TestStatus) Int32() int32   { return int32(s) }

var TestStatuses = []TestStatus{TestStopped, TestRunning, TestExited}

func TestState(t *testing.T) {
	s, err := NewState(TestStopped, TestRunning, TestExited)
	if err != nil {
		t.Fatal(err)
	}
	stats := TestStatuses
	for _, st := range stats {
		s.Set(st)
		if !s.Is(st) {
			t.Errorf("Is: invalid state: %v", st)
		}
		if s.Status() != st {
			t.Errorf("Get: invalid state: %v", st)
		}
		if !s.CompareAndSwap(st, st) {
			t.Errorf("CompareAndSwap: invalid state: %v", st)
		}
	}
	for i := 1; i < len(stats); i++ {
		old := stats[i-1]
		new := stats[i]
		s.Set(old)
		if !s.CompareAndSwap(old, new) {
			t.Errorf("CompareAndSwap: invalid state: %v", new)
		}
		if s.Is(old) {
			t.Errorf("Is: invalid state: %v", new)
		}
	}
}

func BenchmarkIs(b *testing.B) {
	var states = [...]TestStatus{
		TestStopped,
		TestRunning,
		TestExited,
	}

	done := make(chan struct{})
	defer close(done)
	s, err := NewState(TestStopped, TestRunning, TestExited)
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		n := 0
		for {
			select {
			case <-done:
				return
			default:
				s.Set(states[n&int(len(states)-1)])
			}
			n++
		}
	}()
	for i := 0; i < b.N; i++ {
		s.Is(states[i&int(len(states)-1)])
	}
}
