package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	g   *errgroup.Group
	ctx context.Context
}

func New(ctx context.Context, funcs ...func(context.Context) error) (g *Group) {
	g = &Group{}
	g.g, g.ctx = errgroup.WithContext(ctx)
	g.Go(funcs...)
	return g
}

func (g *Group) Go(funcs ...func(context.Context) error) {
	for _, f := range funcs {
		fn := f
		g.g.Go(func() error {
			return fn(g.ctx)
		})
	}
}

func (g *Group) Wait() error {
	return g.g.Wait()
}
