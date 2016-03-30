package dispatcher

import (
	"encoding/json"
	"reflect"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	bgcaction "bosh-google-cpi/action"
)

type Caller interface {
	Call(bgcaction.Action, []interface{}) (interface{}, error)
}

// JSONCaller unmarshals call arguments with json package and calls action.Run
type JSONCaller struct{}

func NewJSONCaller() JSONCaller {
	return JSONCaller{}
}

func (r JSONCaller) Call(action bgcaction.Action, args []interface{}) (value interface{}, err error) {
	actionValue := reflect.ValueOf(action)
	runMethodValue := actionValue.MethodByName("Run")
	if runMethodValue.Kind() != reflect.Func {
		err = bosherr.Error("Run method not found")
		return
	}

	runMethodType := runMethodValue.Type()
	if r.invalidReturnTypes(runMethodType) {
		err = bosherr.Error("Run method should return a value and an error")
		return
	}

	methodArgs, err := r.extractMethodArgs(runMethodType, args)
	if err != nil {
		err = bosherr.WrapError(err, "Extracting method arguments from payload")
		return
	}

	values := runMethodValue.Call(methodArgs)
	return r.extractReturns(values)
}

func (r JSONCaller) invalidReturnTypes(methodType reflect.Type) (valid bool) {
	if methodType.NumOut() != 2 {
		return true
	}

	secondReturnType := methodType.Out(1)
	if secondReturnType.Kind() != reflect.Interface {
		return true
	}

	errorType := reflect.TypeOf(bosherr.Error(""))
	secondReturnIsError := errorType.Implements(secondReturnType)
	if !secondReturnIsError {
		return true
	}

	return
}

func (r JSONCaller) extractMethodArgs(runMethodType reflect.Type, args []interface{}) (methodArgs []reflect.Value, err error) {
	numberOfArgs := runMethodType.NumIn()
	numberOfReqArgs := numberOfArgs

	if runMethodType.IsVariadic() {
		numberOfReqArgs--
	}

	if len(args) < numberOfReqArgs {
		err = bosherr.Errorf("Not enough arguments, expected %d, got %d", numberOfReqArgs, len(args))
		return
	}

	for i, argFromPayload := range args {
		var rawArgBytes []byte
		rawArgBytes, err = json.Marshal(argFromPayload)
		if err != nil {
			err = bosherr.WrapError(err, "Marshalling action argument")
			return
		}

		argType, typeFound := r.getMethodArgType(runMethodType, i)
		if !typeFound {
			continue
		}

		argValuePtr := reflect.New(argType)

		if err = json.Unmarshal(rawArgBytes, argValuePtr.Interface()); err != nil {
			err = bosherr.WrapError(err, "Unmarshalling action argument")
			return
		}

		methodArgs = append(methodArgs, reflect.Indirect(argValuePtr))
	}

	return
}

func (r JSONCaller) getMethodArgType(methodType reflect.Type, index int) (argType reflect.Type, found bool) {
	numberOfArgs := methodType.NumIn()

	switch {
	case !methodType.IsVariadic() && index >= numberOfArgs:
		return nil, false

	case methodType.IsVariadic() && index >= numberOfArgs-1:
		sliceType := methodType.In(numberOfArgs - 1)
		return sliceType.Elem(), true

	default:
		return methodType.In(index), true
	}
}

func (r JSONCaller) extractReturns(values []reflect.Value) (value interface{}, err error) {
	errValue := values[1]
	if !errValue.IsNil() {
		err = errValue.Interface().(error)
	}

	value = values[0].Interface()
	return
}
