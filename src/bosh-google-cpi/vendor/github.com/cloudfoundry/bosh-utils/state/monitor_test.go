package state

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var _ = runtime.GOMAXPROCS(runtime.NumCPU())

func TestMonitorExit(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	var val int64
	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		p := c.NewProcess()
		wg.Add(1)
		go func(p *Process, addr *int64) {
			defer wg.Done()
			for {
				p.Wait()
				atomic.AddInt64(addr, 1)
			}
		}(p, &val)
	}
	c.Start()
	time.Sleep(time.Millisecond * 10)
	c.Exit()
	if c.Status() != Stopped {
		t.Errorf("TestMonitorExit: invalid status: %s", c.Status())
	}
	ch := make(chan struct{})
	go func() {
		c.NewProcess().Wait()
		ch <- struct{}{}
	}()
	c.Start()
	const timeout = time.Millisecond * 50
	select {
	case <-ch:
		// Ok
	case <-time.After(timeout):
		t.Fatalf("TestMonitorExit: timed out after: %s", timeout)
	}
}

func TestMonitor(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// Ensure that new Processes may be added after a
	// call to Exit.
	for i := 0; i < 5; i++ {
		var val int64
		wg := new(sync.WaitGroup)
		for i := 0; i < 100; i++ {
			p := c.NewProcess()
			wg.Add(1)
			go func(p *Process, addr *int64) {
				defer wg.Done()
				for {
					p.Wait()
					atomic.AddInt64(addr, 1)
				}
			}(p, &val)
		}
		c.Start()
		time.Sleep(time.Millisecond)
		if atomic.LoadInt64(&val) == 0 {
			t.Error("Start: expected positive value")
		}

		c.Stop()
		time.Sleep(time.Millisecond)
		atomic.StoreInt64(&val, 0)
		if n := atomic.LoadInt64(&val); n != 0 {
			t.Errorf("Stop: expected value to equal 0 got: %d", n)
		}

		c.Start()
		time.Sleep(time.Millisecond)
		if atomic.LoadInt64(&val) == 0 {
			t.Error("Start: expected positive value")
		}

		c.Exit()
		ch := make(chan struct{}, 1)
		go func() {
			wg.Wait()
			ch <- struct{}{}
		}()
		const timeout = time.Millisecond * 50
		select {
		case <-ch:
			// Ok
		case <-time.After(timeout):
			t.Fatalf("Exit: timed out after: %s", timeout)
		}
	}
}
