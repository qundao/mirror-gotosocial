package runners

import (
	"fmt"
	"unsafe"

	"sync"
)

// Processor acts similarly to a sync.Once object, except that it is reusable. After
// the first call to Process(), any further calls before this first has returned will
// block until the first call has returned, and return the same error. This ensures
// that only a single instance of it is ever running at any one time.
type Processor struct{ p atomic_pointer }

// Process will process the given function if first-call, else blocking until
// the first function has returned, returning the same error result.
func (p *Processor) Process(proc func() error) (err error) {
	var i *proc_instance

	for {
		// Attempt to load existing instance.
		ptr := (*proc_instance)(p.p.Load())
		if ptr != nil {

			// Wait on existing.
			ptr.wait.Wait()
			err = ptr.err
			return
		}

		if i == nil {
			// Allocate instance.
			i = new(proc_instance)
			i.wait.Add(1)
		}

		// Try to acquire start slot by
		// setting ptr to *our* instance.
		if p.p.CAS(nil, unsafe.Pointer(i)) {
			defer func() {
				if r := recover(); r != nil {
					if i.err != nil {
						rOld := r // wrap the panic so we don't lose existing returned error
						r = fmt.Errorf("panic occured after error %q: %v", i.err.Error(), rOld)
					}

					// Catch panics and wrap as error return.
					i.err = fmt.Errorf("caught panic: %v", r)
				}

				// Set return.
				err = i.err

				// Release the
				// goroutines.
				i.wait.Done()

				// Free processor.
				p.p.Store(nil)
			}()

			// Run func.
			i.err = proc()
			return
		}
	}
}

type proc_instance struct {
	wait sync.WaitGroup
	err  error
}
