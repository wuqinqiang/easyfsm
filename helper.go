package easyfsm

import (
	"context"
	"github.com/wuqinqiang/easyfsm/log"
)

func GoSafe(fn func()) {
	go goSafe(fn)
}

func goSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger.Log(log.LevelDebug, "defer err:", err)
		}
	}()
	fn()
}

type ForkCtxInterface interface {
	ForkCtx(ctx context.Context) context.Context
}
