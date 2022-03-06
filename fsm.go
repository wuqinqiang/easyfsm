package easyfsm

import (
	"github.com/wuqinqiang/easyfsm/log"
)

// FSM is the finite state machine
type FSM struct {
	// current state
	state State
	// current business
	businessName BusinessName
}

// NewFSM creates a new fsm
func NewFSM(businessName BusinessName, initState State) (fsm *FSM) {
	fsm = new(FSM)
	fsm.state = initState
	fsm.businessName = businessName
	return
}

// Call call the state's event func
func (f *FSM) Call(eventName EventName, opts ...ParamOption) (State, error) {
	businessMap, ok := stateMachineMap[f.businessName]
	if !ok {
		return f.getState(), UnKnownBusinessError{businessName: f.businessName}
	}
	events, ok := businessMap[f.getState()]
	if !ok || events == nil {
		return f.getState(), UnKnownStateError{businessName: f.businessName, state: f.getState()}
	}

	opt := new(Param)
	for _, fn := range opts {
		fn(opt)
	}

	eventEntity, ok := events[eventName]
	if !ok || eventEntity == nil {
		return f.getState(), UnKnownEventError{businessName: f.businessName, state: f.getState(), event: eventName}
	}

	// call eventName func
	state, err := eventEntity.Execute(opt)
	if err != nil {
		return f.getState(), err
	}
	oldState := f.getState()
	f.setState(state)
	log.DefaultLogger.Log(log.LevelInfo, "eventName:", eventName,
		"beforeState:", oldState, "afterState:", f.getState())
	return f.getState(), nil
}

// getState get the state
func (f *FSM) getState() State {
	return f.state
}

// setState set the state
func (f *FSM) setState(newState State) {
	f.state = newState
}
