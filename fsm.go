package easyfsm

import (
	"fmt"
	"github.com/wuqinqiang/easyfsm/log"
)

// FSM 有限状态机
type FSM struct {
	state        State        // 当前状态
	businessName BusinessName //业务归属
}

// NewFSM 实例化 FSM
func NewFSM(businessName BusinessName, initState State) (fsm *FSM) {
	fsm = new(FSM)
	fsm.state = initState
	fsm.businessName = businessName
	return
}

// 获取当前状态
func (f *FSM) getState() State {
	return f.state
}

// 设置当前状态
func (f *FSM) setState(newState State) {
	f.state = newState
}

// Call 事件处理
func (f *FSM) Call(event EventName, opts ...ParamOption) (State, error) {
	stateMap, ok := stateEventEntityMap[f.businessName]
	if !ok {
		return 0, fmt.Errorf("[警告] business_name:%v 没有注册", f.businessName)
	}
	events, ok := stateMap[f.getState()]
	if !ok || events == nil {
		return 0, fmt.Errorf("[警告] 状态(%v)未定义任何事件", f.getState())
	}

	opt := new(Param)
	for _, fn := range opts {
		fn(opt)
	}

	eventEntity, ok := events[event]
	if !ok || eventEntity == nil {
		return 0, fmt.Errorf("[警告] 状态(%v)不允许操作(%v)", f.getState(), event)
	}

	state, err := eventEntity.handlerEvent(opt)
	if err != nil {
		return 0, err
	}
	oldState := f.getState()
	f.setState(state)
	newState := f.getState()
	log.DefaultLogger.Log(log.LevelInfo, "操作:", event, "状态从:", oldState, "变成:", newState)
	return f.getState(), nil
}
