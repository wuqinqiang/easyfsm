Translate to: [简体中文](https://github.com/wuqinqiang/easyfsm/blob/main/README_zh-CN.md)

### About Easyfsm

a super easy to use finite state machine implemented in go.

its has the following features:
- Easy to use and quick to understand
- Only one global register is required,no need to register in multiple places 
- Support different business->same state value->customize different event handler

Design:

![easyfsm](https://cdn.syst.top/easyfsm2.png)

why we need to differentiate our business?

Because most of the business state values from database,such as the order state is in the order table and the product state in the product table,it is possible that the values are the same.

the same business corresponding to the same attribute state value expression single. different business under the attribute state may appear the same value, but the meaning expressed is different.
```go
fsm:=NewFsm("businessName","currentState")
currentState,err:=fsm.Call("eventName","eventParam")
```
Explain：
- Business:For example there are product state business,order state business.....
- State：to be paid , to be shipped....
- Event：The set of state reachable events. For example, the only reachable events for the pending payment state are: payment events and cancellation events (depending on your business)
- Execution event subject：Execute custom event functions, if necessary, you can also customize the execution of events before and after hook, event subscribers (such as when the  payment events occurs, asynchronous notification of users, etc.)

### UseAge

```go
go get -u  github.com/wuqinqiang/easyfsm
```

Example，

```go
package main

import (
	"fmt"
	"github.com/wuqinqiang/easyfsm"
)

var (
	// business
	businessName easyfsm.BusinessName = "order"

	// states
	initState easyfsm.State = 1 // Initialization
	paidState easyfsm.State = 2 // Paid
	canceled  easyfsm.State = 3 // Canceled

	//events
	paymentOrderEventName easyfsm.EventName = "paymentOrderEventName"
	cancelOrderEventName  easyfsm.EventName = "cancelOrderEventName"
)

type (
	orderParam struct {
		OrderNo string
	}
)

func init() {
	// Payment order event
	entity := easyfsm.NewEventEntity(paymentOrderEventName,
		func(opt *easyfsm.Param) (easyfsm.State, error) {
			param, ok := opt.Data.(orderParam)
			if !ok {
				panic("param err")
			}
			fmt.Printf("param:%+v\n", param)
			// Handling core business
			return paidState, nil
		})

	// Cancellation Event
	cancelEntity := easyfsm.NewEventEntity(cancelOrderEventName,
		func(opt *easyfsm.Param) (easyfsm.State, error) {
			// Handling core business
			param, ok := opt.Data.(orderParam)
			if !ok {
				panic("param err")
			}
			fmt.Printf("param:%+v\n", param)
			return canceled, nil
		})

	// Register Order State Machine
	easyfsm.RegisterStateMachine(businessName,
		initState,
		entity, cancelEntity)
}

func main() {

	// Normal operation

	// The first step generates the fsm based on the business, and the current state
	fsm := easyfsm.NewFSM(businessName, initState)

	// Step 2:Call the event
	currentState, err := fsm.Call(cancelOrderEventName,
		easyfsm.WithData(orderParam{OrderNo: "wuqinqiang050@gmail.com"}))

	fmt.Printf("[Success]call cancelOrderEventName err:%v\n", err)
	fmt.Printf("[Success]call cancelOrderEventName state:%v\n", currentState)

	// Exception 1, no goods business defined
	fsm = easyfsm.NewFSM("goods", paidState)
	currentState, err = fsm.Call(cancelOrderEventName,
		easyfsm.WithData(orderParam{OrderNo: "wuqinqiang050@gmail.com"}))
	fmt.Printf("[UnKnowBusiness]faild: %v\n", err)
	fmt.Printf("[UnKnowBusiness]faild state:%v\n", currentState)

	//Exception 2, no state defined:2
	fsm = easyfsm.NewFSM(businessName, easyfsm.State(2))
	currentState, err = fsm.Call(cancelOrderEventName,
		easyfsm.WithData(orderParam{OrderNo: "wuqinqiang050@gmail.com"}))
	fmt.Printf("[UnKnowState]faild: %v\n", err)
	fmt.Printf("[UnKnowState]faild state:%v\n", currentState)

	// Exception 3,The shipping event corresponding to state 1 is not defined
	fsm = easyfsm.NewFSM(businessName, initState)
	currentState, err = fsm.Call("shippingEvent",
		easyfsm.WithData(orderParam{OrderNo: "wuqinqiang050@gmail.com"}))
	fmt.Printf("[UnKnowEvent]faild: %v\n", err)
	fmt.Printf("[UnKnowEvent]faild state:%v\n", currentState)

}
```



### Hook

If you want to execute some hooks before and after the event handling function, or execute some other operations asynchronously after the event execution, easyfsm defines these two interfaces.

```go
type (
	EventObserver interface {
		Receive(opt *Param)
	}

	EventHook interface {
		Before(opt *Param)
		After(opt Param, state State, err error)
	}
)
```

We can implement these two interfaces that

```go
type (
	NotifyExample struct {
	}
	HookExample struct {
	}
)

func (h HookExample) Before(opt *easyfsm.Param) {
     fmt.Println("Before event execution")
}

func (h HookExample) After(opt easyfsm.Param, state easyfsm.State, err error) {
     fmt.Println("After event execution")
}

func (o NotifyExample) Receive(opt *easyfsm.Param) {
     fmt.Println("Receive events, send messages")
}
```

Full example code：

```go
package main

import (
	"fmt"
	"time"

	"github.com/wuqinqiang/easyfsm"
)

var (
	// business
	businessName easyfsm.BusinessName = "order"

	// states
	initState easyfsm.State = 1 // initial state
	paidState easyfsm.State = 2 // paid
	canceled  easyfsm.State = 3 // canceled

	//events
	paymentOrderEventName easyfsm.EventName = "paymentOrderEventName"
	cancelOrderEventName  easyfsm.EventName = "cancelOrderEventName"
)

type (
	orderParam struct {
		OrderNo string
	}
)

var (
	_ easyfsm.EventObserver = (*NotifyExample)(nil)
	_ easyfsm.EventHook     = (*HookExample)(nil)
)

type (
	NotifyExample struct {
	}
	HookExample struct {
	}
)

func (h HookExample) Before(opt *easyfsm.Param) {
	fmt.Println("Before event execution")
}

func (h HookExample) After(opt easyfsm.Param, state easyfsm.State, err error) {
	fmt.Println("After event execution")
}

func (o NotifyExample) Receive(opt *easyfsm.Param) {
	fmt.Println("Receive events, send messages")
}

func init() {
	// Payment order event
	entity := easyfsm.NewEventEntity(paymentOrderEventName,
		func(opt *easyfsm.Param) (easyfsm.State, error) {
			param, ok := opt.Data.(orderParam)
			if !ok {
				panic("param err")
			}
			fmt.Printf("param:%+v\n", param)
			// Handling core business
			return paidState, nil
		}, easyfsm.WithHook(HookExample{}), easyfsm.WithObservers(NotifyExample{}))

	// Cancellation Event
	cancelEntity := easyfsm.NewEventEntity(cancelOrderEventName,
		func(opt *easyfsm.Param) (easyfsm.State, error) {
			// Handling core business
			param, ok := opt.Data.(orderParam)
			if !ok {
				panic("param err")
			}
			fmt.Printf("param:%+v\n", param)
			return canceled, nil
		}, easyfsm.WithHook(HookExample{}))

	// Register Order State Machine
	easyfsm.RegisterStateMachine(businessName,
		initState,
		entity, cancelEntity)
}

func main() {

	// Normal operation

	// The first step generates the fsm based on the business, and the current state
	fsm := easyfsm.NewFSM(businessName, initState)

	// Step 2:Call the event
	currentState, err := fsm.Call(paymentOrderEventName,
		easyfsm.WithData(orderParam{OrderNo: "wuqinqiang050@gmail.com"}))

	fmt.Printf("[Success]call paymentOrderEventName err:%v\n", err)
	fmt.Printf("[Success]call paymentOrderEventName state:%v\n", currentState)
	time.Sleep(2 * time.Second)
}
```



### End

If there are different needs, you are welcome to leave a message in issue.
