package easyfsm

import (
	"context"
)

type (
	ParamOption func(*Param)
)

// Param is a the parameter you need to execute the event function.
type Param struct {
	ctx   context.Context
	param interface{}
}

func WithCtx(ctx context.Context) ParamOption {
	return func(opt *Param) {
		opt.ctx = ctx
	}
}

func WithParam(param interface{}) ParamOption {
	return func(opt *Param) {
		opt.param = param
	}
}
