package errorsForProject

import "github.com/pkg/errors"

type RegisterError struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

var (
	ErrCantMarshal = errors.New("cant Marshal")
	ErrCantDelete  = errors.New("cant Delete")
)
