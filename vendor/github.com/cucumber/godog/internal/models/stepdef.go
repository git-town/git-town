package models

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	messages "github.com/cucumber/messages/go/v21"

	"github.com/cucumber/godog/formatters"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

// matchable errors
var (
	ErrUnmatchedStepArgumentNumber = errors.New("func expected more arguments than given")
	ErrCannotConvert               = errors.New("cannot convert argument")
	ErrUnsupportedParameterType    = errors.New("func has unsupported parameter type")
)

// StepDefinition ...
type StepDefinition struct {
	formatters.StepDefinition

	Args         []interface{}
	HandlerValue reflect.Value

	// multistep related
	Nested    bool
	Undefined []string
}

var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

// Run a step with the matched arguments using reflect
// Returns one of ...
// (context, error)
// (context, godog.Steps)
func (sd *StepDefinition) Run(ctx context.Context) (context.Context, interface{}) {
	var values []reflect.Value

	typ := sd.HandlerValue.Type()
	numIn := typ.NumIn()
	hasCtxIn := numIn > 0 && typ.In(0).Implements(typeOfContext)
	ctxOffset := 0

	if hasCtxIn {
		values = append(values, reflect.ValueOf(ctx))
		ctxOffset = 1
		numIn--
	}

	if len(sd.Args) < numIn {
		return ctx, fmt.Errorf("%w: expected %d arguments, matched %d from step", ErrUnmatchedStepArgumentNumber, typ.NumIn(), len(sd.Args))
	}

	for i := 0; i < numIn; i++ {
		param := typ.In(i + ctxOffset)
		switch param.Kind() {
		case reflect.Int:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to int: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int(v)))
		case reflect.Int64:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to int64: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(v))
		case reflect.Int32:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to int32: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int32(v)))
		case reflect.Int16:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to int16: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int16(v)))
		case reflect.Int8:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to int8: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int8(v)))
		case reflect.String:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			values = append(values, reflect.ValueOf(s))
		case reflect.Float64:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to float64: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(v))
		case reflect.Float32:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return ctx, err
			}
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return ctx, fmt.Errorf(`%w %d: "%s" to float32: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(float32(v)))
		case reflect.Ptr:
			arg := sd.Args[i]
			switch param.Elem().String() {
			case "messages.PickleDocString":
				if v, ok := arg.(*messages.PickleStepArgument); ok {
					values = append(values, reflect.ValueOf(v.DocString))
					break
				}

				if v, ok := arg.(*messages.PickleDocString); ok {
					values = append(values, reflect.ValueOf(v))
					break
				}

				return ctx, fmt.Errorf(`%w %d: "%v" of type "%T" to *messages.PickleDocString`, ErrCannotConvert, i, arg, arg)
			case "messages.PickleTable":
				if v, ok := arg.(*messages.PickleStepArgument); ok {
					values = append(values, reflect.ValueOf(v.DataTable))
					break
				}

				if v, ok := arg.(*messages.PickleTable); ok {
					values = append(values, reflect.ValueOf(v))
					break
				}

				return ctx, fmt.Errorf(`%w %d: "%v" of type "%T" to *messages.PickleTable`, ErrCannotConvert, i, arg, arg)
			default:
				// the error here is that the declared function has an unsupported param type - really this ought to be trapped at registration ti,e
				return ctx, fmt.Errorf("%w: the data type of parameter %d type *%s is not supported", ErrUnsupportedParameterType, i, param.Elem().String())
			}
		case reflect.Slice:
			switch param {
			case typeOfBytes:
				s, err := sd.shouldBeString(i)
				if err != nil {
					return ctx, err
				}
				values = append(values, reflect.ValueOf([]byte(s)))
			default:
				// the problem is the function decl is not using a support slice type as the param
				return ctx, fmt.Errorf("%w: the slice parameter %d type []%s is not supported", ErrUnsupportedParameterType, i, param.Elem().Kind())
			}
		case reflect.Struct:
			return ctx, fmt.Errorf("%w: the struct parameter %d type %s is not supported", ErrUnsupportedParameterType, i, param.String())
		default:
			return ctx, fmt.Errorf("%w: the parameter %d type %s is not supported", ErrUnsupportedParameterType, i, param.Kind())
		}
	}

	res := sd.HandlerValue.Call(values)
	if len(res) == 0 {
		return ctx, nil
	}

	// Note that the step fn return types were validated at Initialise in test_context.go stepWithKeyword()

	// single return value may be one of ...
	// error
	// context.Context
	// godog.Steps
	result0 := res[0].Interface()
	if len(res) == 1 {

		// if the single return value is a context then just return it
		if ctx, ok := result0.(context.Context); ok {
			return ctx, nil
		}

		// return type is presumably one of nil, "error" or "Steps" so place it into second return position
		return ctx, result0
	}

	// multi-value value return must be
	//  (context, error) and the context value must not be nil
	if ctx, ok := result0.(context.Context); ok {
		return ctx, res[1].Interface()
	}

	result1 := res[1].Interface()
	errMsg := ""
	if result1 != nil {
		errMsg = fmt.Sprintf(", step def also returned an error: %v", result1)
	}

	text := sd.StepDefinition.Expr.String()

	if result0 == nil {
		panic(fmt.Sprintf("step definition '%v' with return type (context.Context, error) must not return <nil> for the context.Context value%s", text, errMsg))
	}

	panic(fmt.Errorf("step definition '%v' has return type (context.Context, error), but found %v rather than a context.Context value%s", text, result0, errMsg))
}

func (sd *StepDefinition) shouldBeString(idx int) (string, error) {
	arg := sd.Args[idx]
	switch arg := arg.(type) {
	case string:
		return arg, nil
	case *messages.PickleStepArgument:
		if arg.DocString == nil {
			return "", fmt.Errorf(`%w %d: "%v" of type "%T": DocString is not set`, ErrCannotConvert, idx, arg, arg)
		}
		return arg.DocString.Content, nil
	case *messages.PickleDocString:
		return arg.Content, nil
	default:
		return "", fmt.Errorf(`%w %d: "%v" of type "%T" to string`, ErrCannotConvert, idx, arg, arg)
	}
}

// GetInternalStepDefinition ...
func (sd *StepDefinition) GetInternalStepDefinition() *formatters.StepDefinition {
	if sd == nil {
		return nil
	}

	return &sd.StepDefinition
}
