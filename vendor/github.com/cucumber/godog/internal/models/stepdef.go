package models

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/formatters"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

// matchable errors
var (
	ErrUnmatchedStepArgumentNumber = errors.New("func received more arguments than expected")
	ErrCannotConvert               = errors.New("cannot convert argument")
	ErrUnsupportedArgumentType     = errors.New("unsupported argument type")
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
				return ctx, fmt.Errorf("%w: the argument %d type %T is not supported %s", ErrUnsupportedArgumentType, i, arg, param.Elem().String())
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
				return ctx, fmt.Errorf("%w: the slice argument %d type %s is not supported", ErrUnsupportedArgumentType, i, param.Kind())
			}
		default:
			return ctx, fmt.Errorf("%w: the argument %d type %s is not supported", ErrUnsupportedArgumentType, i, param.Kind())
		}
	}

	res := sd.HandlerValue.Call(values)
	if len(res) == 0 {
		return ctx, nil
	}

	if len(res) == 1 {
		r := res[0].Interface()

		if ctx, ok := r.(context.Context); ok {
			return ctx, nil
		}

		return ctx, res[0].Interface()
	}

	return res[0].Interface().(context.Context), res[1].Interface()
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
