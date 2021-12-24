package apiserver

import (
	"errors"
	"net"
	"os"
	"syscall"
)

// Listen starts a listener on the given unix socket.
//
// It deletes the socket file if it already exists.
func Listen(socket string) (net.Listener, error) {
	l, err := net.Listen("unix", socket)
	if err == nil {
		return l, nil
	}

	if !errors.Is(err, syscall.EADDRINUSE) {
		return nil, err
	}

	if err := os.Remove(socket); err != nil {
		return nil, err
	}

	return net.Listen("unix", socket)
}
