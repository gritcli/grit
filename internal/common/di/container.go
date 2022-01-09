package di

import (
	"sync"

	"go.uber.org/dig"
	"go.uber.org/multierr"
)

// Container is a dependency injection container.
type Container struct {
	conM sync.Mutex
	con  *dig.Container

	deferredM sync.Mutex
	deferred  []func() error
}

// Provide registers a new provider with the container.
func (c *Container) Provide(fn interface{}) {
	if err := c.container().Provide(fn); err != nil {
		panic(err)
	}
}

// Invoke invokes fn with arguments supplied by the container.
func (c *Container) Invoke(fn interface{}) error {
	err := c.container().Invoke(fn)
	return unwrapError(err)
}

// Defer registers a function to be called when the container is closed.
func (c *Container) Defer(fn func() error) {
	c.deferredM.Lock()
	c.deferred = append(c.deferred, fn)
	c.deferredM.Unlock()
}

// Close closes the container, calling any functions that were deferred via a
// Deferrer.
func (c *Container) Close() error {
	c.deferredM.Lock()
	deferred := c.deferred
	c.deferred = nil
	c.deferredM.Unlock()

	var err error

	for _, fn := range deferred {
		fn := fn // capture loop variable

		// Use defer construct so that we get the same panic-handling semantics
		// as usual.
		defer func() {
			err = multierr.Append(
				err,
				fn(),
			)
		}()
	}

	return err
}

func (c *Container) container() *dig.Container {
	c.conM.Lock()
	defer c.conM.Unlock()

	if c.con == nil {
		c.con = dig.New()
	}

	return c.con
}
