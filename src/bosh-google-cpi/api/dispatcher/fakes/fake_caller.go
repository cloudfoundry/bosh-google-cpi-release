package fakes

import (
	bgcaction "bosh-google-cpi/action"
)

type FakeCaller struct {
	CallAction bgcaction.Action
	CallArgs   []interface{}
	CallResult interface{}
	CallErr    error
}

func (caller *FakeCaller) Call(action bgcaction.Action, args []interface{}) (interface{}, error) {
	caller.CallAction = action
	caller.CallArgs = args
	return caller.CallResult, caller.CallErr
}
