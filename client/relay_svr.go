package client

import "io"

type RelayTcp interface {
	io.Reader
	io.Writer
	io.Closer
}
