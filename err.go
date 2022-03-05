package easyfsm

import "strconv"

type UnKnownBusinessError struct {
	businessName BusinessName
}

func (businessErr UnKnownBusinessError) Error() string {
	return string(businessErr.businessName) + " is not registered"
}

type UnKnownStateError struct {
	businessName BusinessName
	state        State
}

func (StateErr UnKnownStateError) Error() string {
	return string(StateErr.businessName) + "Not included state:" +
		strconv.Itoa(int(StateErr.state))
}

type UnKnownEventError struct {
	businessName BusinessName
	state        State
	event        EventName
}

func (eventErr UnKnownEventError) Error() string {
	return string(eventErr.businessName) + "Not included event:" +
		string(eventErr.event) + " in state:" + strconv.Itoa(int(eventErr.state))
}
