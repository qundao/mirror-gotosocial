package sched

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"codeberg.org/gruf/go-runners"
)

// Precision is the maximum time we can
// offer scheduler run-time precision down to.
const Precision = 2 * time.Millisecond

var (
	// neverticks is a timer channel
	// that never ticks (it's starved).
	neverticks = make(chan time.Time)

	// alwaysticks is a timer channel
	// that always ticks (it's closed).
	alwaysticks = func() chan time.Time {
		ch := make(chan time.Time)
		close(ch)
		return ch
	}()
)

// Scheduler provides a means of running jobs at specific times and
// regular intervals, all while sharing a single underlying timer.
type Scheduler struct {
	svc  runners.Service // svc manages the main scheduler routine
	jobs []*Job          // jobs is a list of tracked Jobs to be executed
	jch  atomic_channel  // jch accepts either Jobs or job IDs to notify new/removed jobs
	jid  atomic.Uint64   // jid is used to iteratively generate unique IDs for jobs
}

// Start will attempt to start the Scheduler. Immediately returns false
// if the Service is already running, and true after completed run.
func (sch *Scheduler) Start() bool {
	var wait sync.WaitGroup

	// Use waiter to synchronize between started
	// goroutine and ourselves, to ensure that
	// we don't return before Scheduler init'd.
	wait.Add(1)

	ok := sch.svc.GoRun(func(ctx context.Context) {
		// Prepare new channel.
		ch := new(channel)
		ch.ctx = ctx.Done()
		ch.ch = make(chan interface{})
		sch.jch.Store(ch)

		// Release
		// start fn
		wait.Done()

		// Main loop
		sch.run(ch)
	})

	if ok {
		// Wait on
		// goroutine
		wait.Wait()
	} else {
		// Release
		wait.Done()
	}

	return ok
}

// Stop will attempt to stop the Scheduler. Immediately returns false
// if not running, and true only after Scheduler is fully stopped.
func (sch *Scheduler) Stop() bool {
	return sch.svc.Stop()
}

// Running will return whether Scheduler is running (i.e. NOT stopped / stopping).
func (sch *Scheduler) Running() bool {
	return sch.svc.Running()
}

// Done returns a channel that's closed when Scheduler.Stop() is called.
func (sch *Scheduler) Done() <-chan struct{} {
	return sch.svc.Done()
}

// Schedule will add provided Job to the Scheduler, returning a cancel function.
func (sch *Scheduler) Schedule(job *Job) (cancel func()) {
	if job == nil {
		panic("nil job")
	}

	// Load job channel.
	ch := sch.jch.Load()
	if ch == nil {
		panic("not running")
	}

	// Calculate next job ID.
	job.id = sch.jid.Add(1)

	// Pass job
	// to channel.
	if !ch.w(job) {
		panic("not running")
	}

	// Return cancel function for job
	return func() { ch.w(job.id) }
}

// run is the main scheduler run routine, which runs for as long as ctx is valid.
func (sch *Scheduler) run(ch *channel) {
	defer ch.close()
	if ch == nil {
		panic("nil channel")
	} else if sch == nil {
		panic("nil scheduler")
	}

	var (
		// now stores the current time, and will only be
		// set when the timer channel is set to be the
		// 'alwaysticks' channel. this allows minimizing
		// the number of calls required to time.Now().
		now time.Time

		// timerset represents whether timer was running
		// for a particular run of the loop. false means
		// that tch == neverticks || tch == alwaysticks.
		timerset bool

		// timer tick channel (or always / never ticks).
		tch <-chan time.Time

		// timer notifies this main routine to wake when
		// the job queued needs to be checked for executions.
		timer *time.Timer

		// stopdrain will stop and drain the timer
		// if it has been running (i.e. timerset == true).
		stopdrain = func() {
			if timerset && !timer.Stop() {
				<-timer.C
			}
		}
	)

	// Create a stopped timer.
	timer = time.NewTimer(1)
	<-timer.C

	for {
		// Reset timer state.
		timerset = false

		if len(sch.jobs) > 0 {
			// Get now time.
			now = time.Now()

			// Sort by next occurring.
			sort.Sort(byNext(sch.jobs))

			// Get next job time.
			next := sch.jobs[0].Next()

			// If this job is *just* about to be ready, we don't bother
			// sleeping. It's wasted cycles only sleeping for some obscenely
			// tiny amount of time we can't guarantee precision for.
			if until := next.Sub(now); until <= Precision/1e3 {

				// This job is behind,
				// set to always tick.
				tch = alwaysticks
			} else {

				// Reset timer to period.
				timer.Reset(until)
				timerset = true
				tch = timer.C
			}
		} else {

			// Unset timer
			tch = neverticks
		}

		select {
		// Scheduler stopped
		case <-ch.done():
			stopdrain()
			return

		// Timer ticked,
		// run scheduled.
		case t, ok := <-tch:
			if !ok {
				// 'alwaysticks' returns zero
				// times, BUT 'now' will have
				// been set during above sort.
				t = now
			}
			sch.schedule(t)

		// Received update,
		// handle job/id.
		case v := <-ch.r():
			sch.handle(v)
			stopdrain()
		}
	}
}

// handle takes an interfaces received from Scheduler.jch and handles either:
// - Job --> new job to add.
// - uint64 --> job ID to remove.
func (sch *Scheduler) handle(v interface{}) {
	switch v := v.(type) {
	// New job added
	case *Job:
		// Get current time.
		now := time.Now()

		// Update next call time.
		next := v.timing.Next(now)
		v.next.Store(next)

		// Append this job to queued/
		sch.jobs = append(sch.jobs, v)

	// Job removed
	case uint64:
		for i := 0; i < len(sch.jobs); i++ {
			if sch.jobs[i].id == v {
				// This is the job we're looking for! Drop this.
				sch.jobs = append(sch.jobs[:i], sch.jobs[i+1:]...)
				return
			}
		}
	}
}

// schedule will iterate through the scheduler jobs and
// execute those necessary, updating their next call time.
func (sch *Scheduler) schedule(now time.Time) {
	for i := 0; i < len(sch.jobs); {
		// Scope our own var.
		job := sch.jobs[i]

		// We know these jobs are ordered by .Next(), so as soon
		// as we reach one with .Next() after now, we can return.
		if job.Next().After(now) {
			return
		}

		// Run the job.
		go job.Run(now)

		// Update the next call time.
		next := job.timing.Next(now)
		job.next.Store(next)

		if next.IsZero() {
			// Zero time, this job is done and can be dropped.
			sch.jobs = append(sch.jobs[:i], sch.jobs[i+1:]...)
			continue
		}

		// Iter
		i++
	}
}

// byNext is an implementation of sort.Interface
// to sort Jobs by their .Next() time.
type byNext []*Job

func (by byNext) Len() int {
	return len(by)
}

func (by byNext) Less(i int, j int) bool {
	return by[i].Next().Before(by[j].Next())
}

func (by byNext) Swap(i int, j int) {
	by[i], by[j] = by[j], by[i]
}

// atomic_channel wraps a *channel{} with atomic store / load.
type atomic_channel struct{ p unsafe.Pointer }

func (c *atomic_channel) Load() *channel {
	if p := atomic.LoadPointer(&c.p); p != nil {
		return (*channel)(p)
	}
	return nil
}

func (c *atomic_channel) Store(v *channel) {
	atomic.StorePointer(&c.p, unsafe.Pointer(v))
}

// channel wraps both a context done
// channel and a generic interface channel
// to support safe writing to an underlying
// channel that correctly fails after close.
type channel struct {
	ctx <-chan struct{}
	ch  chan interface{}
}

// done returns internal context channel.
func (ch *channel) done() <-chan struct{} {
	return ch.ctx
}

// r returns internal channel for read.
func (ch *channel) r() chan interface{} {
	return ch.ch
}

// w writes 'v' to channel, or returns false if closed.
func (ch *channel) w(v interface{}) bool {
	select {
	case <-ch.ctx:
		return false
	case ch.ch <- v:
		return true
	}
}

// close closes underlying channel.
func (ch *channel) close() {
	close(ch.ch)
}
