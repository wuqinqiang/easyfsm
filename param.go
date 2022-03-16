package easyfsm

import (
	"context"
)

type (
	ParamOption func(*Param)
)

// Param is a parameter you need to execute the event function.
type Param struct {
	Ctx  context.Context
	Data interface{}
}

func WithCtx(ctx context.Context) ParamOption {
	return func(opt *Param) {
		if ctx == nil {
			return
		}
		opt.Ctx = ctx
	}
}

func WithData(data interface{}) ParamOption {
	return func(opt *Param) {
		opt.Data = data
	}
}
