package runners

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Service provides a means of tracking a single long-running Service, provided protected state
// changes and preventing multiple instances running. Also providing Service state information.
type Service struct{ p atomic_pointer }

// Run will run the supplied function until completion, using given context to propagate cancel.
// Immediately returns false if the Service is already running, and true after completed run.
func (svc *Service) Run(fn func(context.Context)) (ok bool) {
	var ptr *svc_instance

	// Attempt to start.
	ptr, ok = svc.start()
	if !ok {
		return
	}

	// Run given function.
	defer svc.on_done(ptr)
	fn(CancelCtx(ptr.done))
	return
}

// GoRun will run the supplied function until completion in a goroutine, using given context to
// propagate cancel. Immediately returns boolean indicating success, or that service is already running.
func (svc *Service) GoRun(fn func(context.Context)) (ok bool) {
	var ptr *svc_instance

	// Attempt to start.
	ptr, ok = svc.start()
	if !ok {
		return
	}

	go func() {
		// Run given function.
		defer svc.on_done(ptr)
		fn(CancelCtx(ptr.done))
	}()

	return
}

// RunWait is functionally the same as .Run(), but blocks until the first instance of .Run() returns.
func (svc *Service) RunWait(fn func(context.Context)) (ok bool) {
	var ptr *svc_instance

	// Attempt to start.
	ptr, ok = svc.start()
	if !ok {
		<-ptr.done
		return
	}

	// Run given function.
	defer svc.on_done(ptr)
	fn(CancelCtx(ptr.done))
	return
}

// GoRunWait is functionally the same as .RunWait(), but blocks until the first instance of RunWait() returns.
func (svc *Service) GoRunWait(fn func(context.Context)) (ok bool) {
	var ptr *svc_instance

	// Attempt to start.
	ptr, ok = svc.start()
	if !ok {
		<-ptr.done
		return
	}

	go func() {
		// Run given function.
		defer svc.on_done(ptr)
		fn(CancelCtx(ptr.done))
	}()

	return
}

// Stop will attempt to stop the service, cancelling the running function's context. Immediately
// returns false if not running, and true only after Service is fully stopped.
func (svc *Service) Stop() bool {
	return svc.must_get().stop()
}

// Running returns if Service is running (i.e. NOT fully stopped, but may be *stopping*).
func (svc *Service) Running() bool {
	return svc.must_get().running()
}

// Done returns a channel that's closed when Service.Stop() is called. It is
// the same channel provided to the currently running service function.
func (svc *Service) Done() <-chan struct{} {
	return svc.must_get().done
}

func (svc *Service) start() (*svc_instance, bool) {
	ptr := svc.must_get()
	return ptr, ptr.start()
}

func (svc *Service) on_done(ptr *svc_instance) {
	// Ensure stopped.
	ptr.stop_private()

	// Free service.
	svc.p.Store(nil)
}

func (svc *Service) must_get() *svc_instance {
	var newptr *svc_instance

	for {
		// Try to load existing instance.
		ptr := (*svc_instance)(svc.p.Load())
		if ptr != nil {
			return ptr
		}

		if newptr == nil {
			// Allocate new instance.
			newptr = new(svc_instance)
			newptr.done = make(chan struct{})
		}

		// Attempt to acquire slot by setting our ptr.
		if !svc.p.CAS(nil, unsafe.Pointer(newptr)) {
			continue
		}

		return newptr
	}
}

type svc_instance struct {
	wait  sync.WaitGroup
	done  chan struct{}
	state atomic.Uint32
}

const (
	started_bit  = uint32(1) << 0
	stopping_bit = uint32(1) << 1
	finished_bit = uint32(1) << 2
)

func (i *svc_instance) start() (ok bool) {
	// Acquire start by setting 'started' bit.
	switch old := i.state.Or(started_bit); {

	case old&finished_bit != 0:
		// Already finished.

	case old&started_bit == 0:
		// Successfully started!
		i.wait.Add(1)
		ok = true
	}

	return
}

// NOTE: MAY ONLY BE CALLED BY STARTING GOROUTINE.
func (i *svc_instance) stop_private() {
	// Attempt set both stopping and finished bits.
	old := i.state.Or(stopping_bit | finished_bit)

	// Only if we weren't already
	// stopping do we close channel.
	if old&stopping_bit == 0 {
		close(i.done)
	}

	// Release
	// waiters.
	i.wait.Done()
}

func (i *svc_instance) stop() (ok bool) {
	// Attempt to set the 'stopping' bit.
	switch old := i.state.Or(stopping_bit); {

	case old&finished_bit != 0:
		// Already finished.
		return

	case old&started_bit == 0:
		// This was never started
		// to begin with, just mark
		// as fully finished here.
		_ = i.state.Or(finished_bit)
		return

	case old&stopping_bit == 0:
		// We succesfully stopped
		// instance, close channel.
		close(i.done)
		ok = true
	}

	// Wait on stop.
	i.wait.Wait()
	return
}

// running returns whether service was started and
// is not yet finished. that indicates that it may
// have been started and not yet stopped, or that
// it was started, stopped and not yet returned.
func (i *svc_instance) running() bool {
	val := i.state.Load()
	return val&started_bit != 0 &&
		val&finished_bit == 0
}
