package easyfsm

import "context"

type (
	// EventEntity is the core that wraps the basic Event methods.
	EventEntity struct {
		hook      EventHook
		eventName EventName
		observers []EventObserver
		eventFunc EventFunc
		// issue:https://github.com/wuqinqiang/easyfsm/issues/16
		forkCtxFunc func(ctx context.Context) context.Context
	}

	EventEntityOpt func(entity *EventEntity)

	// EventFunc is the function that will be called when the event is triggered.
	EventFunc func(opt *Param) (State, error)

	// EventObserver is the interface.When the event is processed,
	//  it can notify the observers asynchronously and execute their own business.
	EventObserver interface {
		Receive(opt *Param)
	}

	// EventHook is the interface that user can implement to hook the event func.
	EventHook interface {
		Before(opt *Param)
		After(opt Param, state State, err error)
	}
)

// NewEventEntity creates a new EventEntity.
func NewEventEntity(event EventName, handler EventFunc,
	opts ...EventEntityOpt) *EventEntity {
	entity := &EventEntity{
		eventName: event,
		eventFunc: handler,
		observers: make([]EventObserver, 0),
		forkCtxFunc: func(ctx context.Context) context.Context {
			return context.Background()
		},
	}
	for _, opt := range opts {
		opt(entity)
	}
	return entity
}

// WithObservers adds observers to the event.
func WithObservers(observers ...EventObserver) EventEntityOpt {
	return func(entity *EventEntity) {
		if len(observers) == 0 {
			return
		}
		entity.observers = append(entity.observers, observers...)
	}
}

// WithHook adds hook to the event
func WithHook(hook EventHook) EventEntityOpt {
	return func(entity *EventEntity) {
		if hook == nil {
			return
		}
		entity.hook = hook
	}
}

func WithForkCtxFunc(fn func(ctx context.Context) context.Context) EventEntityOpt {
	return func(entity *EventEntity) {
		entity.forkCtxFunc = fn
	}
}

// Execute executes the event.
func (e *EventEntity) Execute(param *Param) (State, error) {
	if e.hook != nil {
		e.hook.Before(param)
	}
	state, err := e.eventFunc(param)
	if e.hook != nil {
		// post operation Not allowed to modify the Data
		forkParam := *param
		e.hook.After(forkParam, state, err)
	}
	if err != nil {
		return state, err
	}

	// Asynchronous notify observers
	GoSafe(func() {
		param.Ctx = e.forkCtxFunc(param.Ctx)
		e.notify(param)
	})
	return state, nil
}

func (e *EventEntity) notify(opt *Param) {
	for _, observer := range e.observers {
		if observer == nil {
			continue
		}
		observer.Receive(opt)
	}
}
