package requests

import (
	"net/http"
)

type doResponse int

const (
	doOK doResponse = iota
	doConnect
	doValidate
	doHandle
)

func do(cl *http.Client, req *http.Request, validators []ResponseHandler, h ResponseHandler) (doResponse, error) {
	res, err := cl.Do(req)
	if err != nil {
		return doConnect, err
	}
	defer res.Body.Close()

	for _, v := range validators {
		if v == nil {
			continue
		}
		if err = v(res); err != nil {
			return doValidate, err
		}
	}
	if err = h(res); err != nil {
		return doHandle, err
	}

	return doOK, nil
}
