package easyfsm

import "sync"

var (
	// 业务归属->状态->事件->事件处理器
	stateEventEntityMap map[BusinessName]map[State]map[EventName]*EventEntity
	locker              sync.Mutex
)

func init() {
	stateEventEntityMap = make(map[BusinessName]map[State]map[EventName]*EventEntity)
}

func RegisterBusinessName(name BusinessName, state State,
	event EventName, entity *EventEntity) {
	locker.Lock()
	defer locker.Unlock()
	if entity == nil {
		return
	}
	if stateEventEntityMap[name] == nil {
		stateEventEntityMap[name] = make(map[State]map[EventName]*EventEntity)
	}
	if stateEventEntityMap[name][state] == nil {
		stateEventEntityMap[name][state] = make(map[EventName]*EventEntity)
	}
	stateEventEntityMap[name][state][event] = entity
}
