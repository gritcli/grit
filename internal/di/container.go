package di

import (
	"sync"

	"go.uber.org/dig"
	"go.uber.org/multierr"
)

// Container is a dependency injection container.
type Container struct {
	c *dig.Container

	m        sync.Mutex
	deferred []func() error
}

// New returns a new dependency injection container.
func New() *Container {
	c := &Container{
		c: dig.New(),
	}

	c.Provide(func() *Container {
		return c
	})

	return c
}

// Provide registers a new provider with the container.
func (c *Container) Provide(fn interface{}) {
	if err := c.c.Provide(fn); err != nil {
		panic(err)
	}
}

// Invoke invokes fn with arguments supplied by the container.
func (c *Container) Invoke(fn interface{}) error {
	err := c.c.Invoke(fn)
	return unwrapError(err)
}

// Defer registers a function to be called when the container is closed.
func (c *Container) Defer(fn func() error) {
	c.m.Lock()
	c.deferred = append(c.deferred, fn)
	c.m.Unlock()
}

// Close closes the container, calling any functions that were deferred via a
// Deferrer.
func (c *Container) Close() error {
	c.m.Lock()
	deferred := c.deferred
	c.deferred = nil
	c.m.Unlock()

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
