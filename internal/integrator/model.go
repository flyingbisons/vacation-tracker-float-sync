package integrator

import "errors"

var (
	ErrorRequestNotFound = errors.New("request not found")
)

type Request struct {
	VtRequestID    string
	FloatTimeoffID uint64
	Created        uint64
}
