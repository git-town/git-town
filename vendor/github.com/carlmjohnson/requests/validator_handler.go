package requests

import (
	"errors"
	"fmt"
	"net/http"
)

// ValidatorHandler composes a Validator and a Handler.
// If the validation check fails, it triggers the handler.
// Any errors from validator or handler will be joined to the error returned.
// If the handler succeeds, the error will matching ErrInvalidHandled.
func ValidatorHandler(v, h ResponseHandler) ResponseHandler {
	return func(res *http.Response) error {
		err1 := v(res)
		if err1 == nil { // passes validation
			return nil
		}
		err2 := h(res)
		if err2 == nil { // successfully handled
			return fmt.Errorf("%w: %w", ErrInvalidHandled, err1)
		}
		return errors.Join(err1, err2)
	}
}

var ErrInvalidHandled = errors.New("handled recovery from invalid response")

// ErrorJSON is a ValidatorHandler that applies DefaultValidator
// and decodes the response as a JSON object
// if the DefaultValidator check fails.
func ErrorJSON(v any) ResponseHandler {
	return ValidatorHandler(DefaultValidator, ToJSON(v))
}
