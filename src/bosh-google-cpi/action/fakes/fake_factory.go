package fakes

import (
	"errors"
	"fmt"

	bgcaction "bosh-google-cpi/action"
)

type FakeFactory struct {
	registeredActions    map[string]*FakeAction
	registeredActionErrs map[string]error
}

func NewFakeFactory() *FakeFactory {
	return &FakeFactory{
		registeredActions:    make(map[string]*FakeAction),
		registeredActionErrs: make(map[string]error),
	}
}

func (f *FakeFactory) Create(method string) (bgcaction.Action, error) {
	if err := f.registeredActionErrs[method]; err != nil {
		return nil, err
	}
	if action := f.registeredActions[method]; action != nil {
		return action, nil
	}
	return nil, errors.New("Action not found")
}

func (f *FakeFactory) RegisterAction(method string, action *FakeAction) {
	if a := f.registeredActions[method]; a != nil {
		panic(fmt.Sprintf("Action is already registered: %v", a))
	}
	f.registeredActions[method] = action
}

func (f *FakeFactory) RegisterActionErr(method string, err error) {
	if e := f.registeredActionErrs[method]; e != nil {
		panic(fmt.Sprintf("Action err is already registered: %v", e))
	}
	f.registeredActionErrs[method] = err
}

type FakeAction struct{}

func (a *FakeAction) Run(payload []byte) (interface{}, error) {
	return nil, nil
}
