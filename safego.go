package easyfsm

import "github.com/wuqinqiang/easyfsm/log"

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
