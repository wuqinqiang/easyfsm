package easyfsm

import "sync"

var (
	// [business] -> [state] -> [event] ->[ eventEntity->handler]
	stateMachineMap map[BusinessName]map[State]map[EventName]*EventEntity
	locker          sync.Mutex
)

func init() {
	stateMachineMap = make(map[BusinessName]map[State]map[EventName]*EventEntity)
}

//RegisterStateMachine register state machine
func RegisterStateMachine(name BusinessName, state State,
	event EventName, entity *EventEntity) {
	locker.Lock()
	defer locker.Unlock()
	if entity == nil {
		return
	}
	if stateMachineMap[name] == nil {
		stateMachineMap[name] = make(map[State]map[EventName]*EventEntity)
	}
	if stateMachineMap[name][state] == nil {
		stateMachineMap[name][state] = make(map[EventName]*EventEntity)
	}
	stateMachineMap[name][state][event] = entity
}
