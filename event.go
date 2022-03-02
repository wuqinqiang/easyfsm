package easyfsm

type (
	EventEntity struct {
		hook      EventHook
		eventName EventName
		observers []EventObserver
		handler   EventFunc
	}

	EventEntityOpt func(entity *EventEntity)

	EventFunc func(opt *Param) (State, error)

	EventObserver interface {
		Receive(opt *Param)
	}

	EventHook interface {
		Before(opt *Param)
		After(opt Param, state State, err error)
	}
)

func NewEventEntity(event EventName, handler EventFunc,
	opts ...EventEntityOpt) *EventEntity {
	entity := &EventEntity{
		eventName: event,
		handler:   handler,
		observers: make([]EventObserver, 0),
	}
	for _, opt := range opts {
		opt(entity)
	}
	return entity
}

func WithObservers(observers ...EventObserver) EventEntityOpt {
	return func(entity *EventEntity) {
		if len(observers) == 0 {
			return
		}
		entity.observers = append(entity.observers, observers...)
	}
}

func WithHook(hook EventHook) EventEntityOpt {
	return func(entity *EventEntity) {
		if hook == nil {
			return
		}
		entity.hook = hook
	}
}

func (e *EventEntity) handlerEvent(param *Param) (State, error) {
	if e.hook != nil {
		e.hook.Before(param)
	}
	state, err := e.handler(param)
	if e.hook != nil {
		forkParam := *param
		e.hook.After(forkParam, state, err)
	}
	if err != nil {
		return state, err
	}

	GoSafe(func() {
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
