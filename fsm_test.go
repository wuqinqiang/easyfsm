package easyfsm

import (
	"reflect"
	"testing"
)

var (
	businessName = BusinessName("business_order")
)

func Init() {
	args := DefaultArgList()
	for i := range DefaultArgList() {
		RegisterStateMachine(businessName, args[i].state,
			args[i].eventName, &args[i].entity)
	}
}

func TestNewFSM(t *testing.T) {
	initState := State(0)
	wantFsm := &FSM{
		state:        initState,
		businessName: businessName,
	}
	fsm := NewFSM(businessName, initState)

	if !reflect.DeepEqual(fsm, wantFsm) {
		t.Errorf("fsm %v want %v", fsm, wantFsm)
	}
}

func TestFSM_Call(t *testing.T) {
	//clear
	stateMachineMap = make(map[BusinessName]map[State]map[EventName]*EventEntity)
	Init()
	type (
		wantRes struct {
			error error
			state State
		}

		arg struct {
			state     State
			eventName EventName
			wantRes   wantRes
		}
	)

	args := []arg{
		{
			state:     State(0),
			eventName: EventName("crateOrder"),
			wantRes: wantRes{
				error: nil,
				state: 1,
			},
		},
		{
			state:     State(1),
			eventName: EventName("payOrder"),
			wantRes: wantRes{
				error: nil,
				state: 2,
			},
		},
		{
			state:     State(1),
			eventName: EventName("cancelOrder"),
			wantRes: wantRes{
				error: nil,
				state: 3,
			},
		},
		{
			state:     State(3),
			eventName: EventName("no_state"),
			wantRes: wantRes{
				error: UnKnownStateError{businessName: businessName, state: State(3)},
				state: State(3),
			},
		},
		{
			state:     State(1),
			eventName: EventName("no_event"),
			wantRes: wantRes{
				error: UnKnownEventError{businessName: businessName,
					event: EventName("no_event"), state: State(1)},
				state: State(1),
			},
		},
	}

	for i := range args {
		fsm := NewFSM(businessName, args[i].state)
		resState, err := fsm.Call(args[i].eventName)
		if err != args[i].wantRes.error {
			t.Errorf("err %v want %v", err, args[i].wantRes.error)
		}
		if resState != args[i].wantRes.state {
			t.Errorf("state %v want %v", resState, args[i].wantRes.state)
		}
	}
}
