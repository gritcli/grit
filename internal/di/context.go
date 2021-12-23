package di

import (
	"context"
)

type containerContextKey struct{}

// ContextWithContainer returns a child of ctx that contains c.
func ContextWithContainer(ctx context.Context, c *Container) context.Context {
	return context.WithValue(ctx, containerContextKey{}, c)
}

// ContainerFromContext extracts a container from ctx, or panics if there is no
// container in ctx.
func ContainerFromContext(ctx context.Context) *Container {
	return ctx.Value(containerContextKey{}).(*Container)
}
