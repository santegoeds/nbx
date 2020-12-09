package errors

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrServer     = errors.New("server error")
	ErrInternal   = errors.New("internal error")
)

func FromResponse(rsp *http.Response, format string, a ...interface{}) error {
	errClass := rsp.StatusCode / 100
	var baseErr error
	switch errClass {
	case 4:
		baseErr = ErrBadRequest
	case 5:
		baseErr = ErrServer
	default:
		baseErr = ErrInternal
	}

	msg := fmt.Sprintf(format, a...)
	serverMsg, _ := ioutil.ReadAll(rsp.Body)
	return fmt.Errorf("%s. %w {status-code: %d, message: %s}", msg, baseErr, rsp.StatusCode, serverMsg)
}
