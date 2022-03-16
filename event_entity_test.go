package easyfsm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var _ EventHook = (*HookTest)(nil)

type HookTest struct {
}

func (h HookTest) Before(opt *Param) {
	//TODO implement me
	panic("implement me")
}

func (h HookTest) After(opt Param, state State, err error) {
	//TODO implement me
	panic("implement me")
}

func funcEqual(a, b interface{}) bool {
	av := reflect.ValueOf(&a).Elem()
	bv := reflect.ValueOf(&b).Elem()
	return av.InterfaceData() == bv.InterfaceData()
}

var _ EventObserver = (*ObserverTest)(nil)

type ObserverTest struct {
}

func (o ObserverTest) Receive(opt *Param) {
	//TODO implement me
	panic("implement me")
}

func TestNewEventEntityNoOption(t *testing.T) {
	businessName := "create with no option"
	eventName := "create_order"
	handler := func(opt *Param) (State, error) {
		return State(1), nil
	}
	testKey := "remember"
	testVal := "ok"

	forkCtxFunc := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, testKey, testVal)
	}

	wantEntity := NewEventEntity(EventName(eventName), handler, WithForkCtxFunc(forkCtxFunc))

	t.Run(businessName, func(t *testing.T) {
		got := NewEventEntity(EventName(eventName), handler, WithForkCtxFunc(forkCtxFunc))

		if !reflect.DeepEqual(got.eventName, wantEntity.eventName) {
			t.Errorf("eventEntity name =%v,want %v", got.eventName, wantEntity.eventName)
		}

		if !reflect.DeepEqual(got.observers, wantEntity.observers) {
			t.Errorf("eventEntity observers =%v,want %v", got.observers, wantEntity.observers)
		}

		if !funcEqual(got.eventFunc, wantEntity.eventFunc) {
			t.Errorf("eventEntity handler =%v,want %v", got.eventFunc, wantEntity.eventFunc)
		}
		ctx := got.forkCtxFunc(context.Background())
		val := ctx.Value(testKey).(string)
		if val != testVal {
			t.Errorf("eventEntity ctxText =%v,want %v", val, testVal)
		}

	})
}

func TestNewEventEntityWithObservers(t *testing.T) {
	businessName := "create with observers"
	eventName := "create_order"
	handler := func(opt *Param) (State, error) {
		return State(1), nil
	}

	wantEntity := NewEventEntity(EventName(eventName), handler, WithObservers(ObserverTest{}))

	t.Run(businessName, func(t *testing.T) {
		got := NewEventEntity(EventName(eventName), handler, WithObservers(ObserverTest{}))
		if len(got.observers) != 1 {
			t.Errorf("eventEntity observers len =%v,want %v", got.observers, wantEntity.observers)
		}
	})
}

func TestNewEventEntityWithHook(t *testing.T) {
	businessName := "create with Hook"
	eventName := "create_order"
	handler := func(opt *Param) (State, error) {
		return State(1), nil
	}

	wantEntity := NewEventEntity(EventName(eventName), handler, WithHook(HookTest{}))

	t.Run(businessName, func(t *testing.T) {
		got := NewEventEntity(EventName(eventName), handler, WithHook(HookTest{}))
		if len(got.observers) != 0 || !reflect.DeepEqual(got.observers, wantEntity.observers) {
			t.Errorf("eventEntity observers =%v,want %v", got.observers, wantEntity.observers)
		}

		if !reflect.DeepEqual(got.hook, wantEntity.hook) {
			t.Errorf("eventEntity hook =%v,want %v", got.hook, wantEntity.hook)
		}

	})
}

func TestEventEntity_Execute_Success(t *testing.T) {
	eventName := "create_order"
	handler := func(opt *Param) (State, error) {
		return State(2), nil
	}
	entity := NewEventEntity(EventName(eventName), handler)

	type CreateOrderPar struct {
		OrderId string
	}
	param := &Param{
		Ctx:  context.TODO(),
		Data: CreateOrderPar{OrderId: "wuqq0223"},
	}
	state, err := entity.execute(param)

	wantState, wantErr := handler(param)
	if err != nil {
		t.Errorf("execute err %v ,want %v", err, wantErr)
	}

	if state != wantState {
		t.Errorf("execute state %v ,want %v", err, wantErr)
	}
}

func TestEventEntity_Execute_Err(t *testing.T) {
	paidErr := fmt.Errorf("paid err")
	eventName := "paid_order"
	handler := func(opt *Param) (State, error) {
		return State(3), paidErr
	}
	entity := NewEventEntity(EventName(eventName), handler)

	type CreateOrderPar struct {
		OrderId string
	}
	param := &Param{
		Ctx:  context.TODO(),
		Data: CreateOrderPar{OrderId: "wuqq0223"},
	}
	state, err := entity.execute(param)

	wantState, wantErr := handler(param)
	if err == nil || !errors.Is(err, paidErr) {
		t.Errorf("execute err %v ,want %v", err, wantErr)
	}
	if state != wantState {
		t.Errorf("execute state %v ,want %v", err, wantErr)
	}
}
