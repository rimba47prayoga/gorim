package middlewares

import "github.com/rimba47prayoga/gorim.git"

type BaseMiddleware struct {
	nextFunc gorim.HandlerFunc
}

func (m *BaseMiddleware) SetNextFunc(next gorim.HandlerFunc) {
	m.nextFunc = next
}

func (m *BaseMiddleware) Next(ctx gorim.Context) error {
	return m.nextFunc(ctx)
}

func (m *BaseMiddleware) Call(c gorim.Context) error {
	panic("Call function must be implemented.")
}

type BaseAuthMiddleware struct {
	BaseMiddleware
}
