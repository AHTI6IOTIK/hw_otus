package telnet

import "errors"

var (
	ErrConnectionFailEstablished = errors.New("connection could not be established")
	ErrConnectionNotEstablished  = errors.New("the connection has not been established")
	ErrConnectionClose           = errors.New("connection closing error")
	ErrConnectionCloseByPear     = errors.New("connection was closed by peer")
)
