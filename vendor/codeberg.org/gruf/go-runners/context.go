package runners

import (
	"context"
	"sync/atomic"
	"time"
)

// closedctx is an always closed context.
var closedctx = func() context.Context {
	ctx := make(chan struct{})
	close(ctx)
	return CancelCtx(ctx)
}()

// Closed returns an always closed context.
func Closed() context.Context {
	return closedctx
}

// CtxWithCancel returns a new context.Context impl with cancel.
func CtxWithCancel() (context.Context, context.CancelFunc) {
	var once atomic.Uint32
	ctx := make(chan struct{})
	return CancelCtx(ctx), func() {
		if once.CompareAndSwap(0, 1) {
			close(ctx)
		}
	}
}

// CancelCtx is the simplest possible cancellable context.
type CancelCtx (<-chan struct{})

func (CancelCtx) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (ctx CancelCtx) Done() <-chan struct{} {
	return ctx
}

func (ctx CancelCtx) Err() error {
	select {
	case <-ctx:
		return context.Canceled
	default:
		return nil
	}
}

func (CancelCtx) Value(key interface{}) interface{} {
	return nil
}

func (ctx CancelCtx) String() string {
	var state string
	select {
	case <-ctx:
		state = "closed"
	default:
		state = "open"
	}
	return "CancelCtx{state:" + state + "}"
}

func (ctx CancelCtx) GoString() string {
	return "runners." + ctx.String()
}
