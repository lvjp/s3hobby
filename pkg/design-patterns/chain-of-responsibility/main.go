package chain_of_responsibility

import (
	"context"
)

type Handler[Request context.Context] interface {
	Handle(request Request) error
}

type HandlerFunc[Request context.Context] func(request Request) error

func (fn HandlerFunc[Request]) Handle(request Request) error {
	return fn(request)
}

type Middleware[Request context.Context] interface {
	Middleware(request Request, next Handler[Request]) error
}

type MiddlewareFunc[Request context.Context] func(request Request, next Handler[Request]) error

func (fn MiddlewareFunc[Request]) Middleware(request Request, next Handler[Request]) error {
	return fn(request, next)
}

type chainLink[Request context.Context] struct {
	next Handler[Request]
	with Middleware[Request]
}

func (chainLink chainLink[Request]) Handle(request Request) error {
	return chainLink.with.Middleware(request, chainLink.next)
}

func NewChain[Request context.Context](h Handler[Request], with ...Middleware[Request]) Handler[Request] {
	for i := len(with) - 1; i >= 0; i-- {
		h = chainLink[Request]{
			next: h,
			with: with[i],
		}
	}

	return h
}
