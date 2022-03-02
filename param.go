package easyfsm

import (
	"context"
)

type (
	ParamOption func(*Param)
)

// Param 定义 EventFunc 所需参数
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
