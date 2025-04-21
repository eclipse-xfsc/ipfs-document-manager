package main

import (
	"fmt"
)

type IpfsError struct {
	Err  error
	Msg  string
	Code int
}

func (e IpfsError) Error() string {
	return fmt.Sprintf("error status: %v: %s with error: %t", e.Code, e.Msg, e.Err)
}
