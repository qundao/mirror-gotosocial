package dereferencing

import (
	"context"

	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/transport"
)

// TODO: at some point i'll move all of the below code
// into gtscontext or in same place we can reduce the
// number of times we have to re-fetch transport types.
// or, alternatively, overhaul the controller type!

type transportContext struct {
	context.Context
	account *gtsmodel.Account
	tsport  transport.Transport
}

func (ctx *transportContext) Value(key any) any {
	requester, ok := key.(*gtsmodel.Account)
	if ok && ctx.account == requester {
		return ctx.tsport
	}
	return ctx.Context.Value(key)
}

// getTransport returns a transport for the requesting account, caching it in the returned context when newly fetched.
func (d *Dereferencer) getTransport(ctx context.Context, requester *gtsmodel.Account) (context.Context, transport.Transport, error) {
	if tsport, ok := ctx.Value(requester).(transport.Transport); ok {
		return ctx, tsport, nil
	}

	// Load new transport instance for the passed requesting account.
	tsport, err := d.transportController.NewTransportForUsername(ctx,
		accountUsername(requester))
	if err != nil {
		return nil, nil, gtserror.Newf("error creating transport for user %s: %w", accountUsername(requester), err)
	}

	// Store the transport in
	// a new wrapping context.
	ctx = &transportContext{
		Context: ctx,
		account: requester,
		tsport:  tsport,
	}

	return ctx, tsport, nil
}
