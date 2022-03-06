package easyfsm

import (
	"context"
)

type (
	ParamOption func(*Param)
)

// Param is a the parameter you need to execute the event function.
type Param struct {
	Ctx  context.Context
	Data interface{}
}

func WithCtx(ctx context.Context) ParamOption {
	return func(opt *Param) {
		opt.Ctx = ctx
	}
}

func WithData(param interface{}) ParamOption {
	return func(opt *Param) {
		opt.Data = param
	}
}
