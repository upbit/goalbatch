// Package goalbatch A simple way to execute functions asynchronously and waits for results
//  The idea come from github.com/vardius/gollback, but more suitable for concurrent waiting scenarios
package goalbatch

import (
	"context"
	"sync"
)

// AsyncFunc represents asynchronous function
type AsyncFunc func(ctx context.Context) (interface{}, error)

// Goalbatch provides set of utility methods to easily manage asynchronous functions
type Goalbatch interface {
	// Batch method returns when all of the callbacks passed or context is done,
	// returned responses and errors are ordered according to callback order
	Batch(fns ...AsyncFunc) ([]interface{}, []error)
}

type goalbatch struct {
	ctx   context.Context
	close bool
	mutex sync.Mutex
}

type response struct {
	idx int
	res interface{}
	err error
}

// Batch method returns when all of the callbacks passed or context is done
// returned responses and errors are ordered according to callback order
func (p *goalbatch) Batch(fns ...AsyncFunc) ([]interface{}, []error) {
	ch := make(chan *response, len(fns))
	defer func() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		p.close = true
		close(ch)
	}()

	for i, fn := range fns {
		go func(index int, f AsyncFunc) {
			var r response
			r.idx = index
			r.res, r.err = f(p.ctx)

			// Check channel if is closed
			p.mutex.Lock()
			defer p.mutex.Unlock()
			if !p.close {
				ch <- &r
			}
		}(i, fn)
	}

	rs := make([]interface{}, len(fns))
	errs := make([]error, len(fns))

	for range fns {
		select {
		case <-p.ctx.Done():
			// context end, may be timeoutd
			return rs, errs
		case r := <-ch:
			index := r.idx
			rs[index] = r.res
			errs[index] = r.err
		}
	}

	return rs, errs
}

// New creates new goalbatch
func New(ctx context.Context) Goalbatch {
	if ctx == nil {
		ctx = context.Background()
	}

	return &goalbatch{
		close: false,
		ctx:   ctx,
	}
}
