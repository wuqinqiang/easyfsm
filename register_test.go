package easyfsm

import (
	"sync"
	"testing"
)

type arg struct {
	state     State
	eventName EventName
	entity    EventEntity
}

func DefaultArgList() []arg {

	var (
		args []arg
	)
	args = append(args, arg{
		state:     0,
		eventName: "crateOrder",
		entity: EventEntity{
			eventName: "crateOrder",
			eventFunc: func(opt *Param) (State, error) {
				return State(1), nil
			},
		},
	},
		arg{
			state:     1,
			eventName: "payOrder",
			entity: EventEntity{
				eventName: "payOrder",
				eventFunc: func(opt *Param) (State, error) {
					return State(2), nil
				},
			},
		},
		arg{
			state:     1,
			eventName: "cancelOrder",
			entity: EventEntity{
				eventName: "cancelOrder",
				eventFunc: func(opt *Param) (State, error) {
					return State(3), nil
				},
			},
		},
	)
	return args
}

func TestRegisterStateMachine(t *testing.T) {
	businessName := BusinessName("business_order")
	args := DefaultArgList()
	// clear
	stateMachineMap = make(map[BusinessName]map[State]map[EventName]*EventEntity)
	for i := range args {
		RegisterStateMachine(businessName, args[i].state,
			args[i].eventName, &args[i].entity)
	}
	commonTest(args, businessName, t)

}

func TestRegisterStateMachineForConcurrent(t *testing.T) {
	businessName := BusinessName("business_order")
	args := DefaultArgList()
	// clear
	stateMachineMap = make(map[BusinessName]map[State]map[EventName]*EventEntity)
	var (
		wg sync.WaitGroup
	)
	for i := range args {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			RegisterStateMachine(businessName, args[index].state,
				args[index].eventName, &args[index].entity)
		}(i)
	}
	wg.Wait()

	commonTest(args, businessName, t)
}

func commonTest(args []arg, businessName BusinessName, t *testing.T) {
	stateMap, ok := stateMachineMap[businessName]
	if !ok {
		t.Errorf("stateMachineMap should have businessName:%v", businessName)
	}
	for j := range args {

		eventMap, ok := stateMap[args[j].state]
		if !ok {
			t.Errorf("stateMachineMap  should have state:%v", args[j].state)
		}
		entity, ok := eventMap[args[j].eventName]
		if !ok {
			t.Errorf("stateMachineMap state %v should have event %v", args[j].state, args[j].eventName)
		}
		if entity == nil {
			t.Errorf("entity  shouldn't be nil")
		}

		state, err := entity.Execute(nil)
		wantState, wantErr := args[j].entity.Execute(nil)
		if err != nil {
			t.Errorf("err %v want:%v", err, wantErr)
		}
		if state != wantState {
			t.Errorf("state %v want:%v", err, wantErr)
		}
	}
}
