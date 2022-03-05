package easyfsm

import (
	"strconv"
	"testing"
)

func TestUnKnownBusinessError_Error(t *testing.T) {
	businessName := BusinessName("order")
	e := UnKnownBusinessError{businessName: businessName}
	if e.Error() != string(businessName)+"is not registered" {
		t.Error("UnKnownBusinessError mismatch")
	}

}

func TestUnKnownEventError_Error(t *testing.T) {
	businessName := BusinessName("order")
	state := State(1)
	event := EventName("pay")
	e := UnKnownEventError{businessName: businessName, state: state, event: event}
	if e.Error() != string(businessName)+"Not included event:"+
		string(event)+" in state:"+strconv.Itoa(int(state)) {
		t.Errorf("UnKnownEventError mismatch")
	}
}

func TestUnKnownStateError_Error(t *testing.T) {
	businessName := BusinessName("order")
	state := State(1)
	e := UnKnownStateError{
		businessName: businessName,
		state:        state,
	}
	if e.Error() != string(businessName)+"Not included state:"+
		strconv.Itoa(int(state)) {
		t.Errorf("UnKnownStateError mismatch")
	}
}
