package apiserver

import (
	"fmt"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/gritcli/grit/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// newStreamLogger returns a logging.Logger that sends log messages over a gRPC
// stream.
func (s *Server) newStreamLogger(
	stream grpc.ServerStream,
	options *api.ClientOptions,
	wrap func(*api.ClientOutput) proto.Message,
) logging.Logger {
	send := func(message string, debug bool) {
		m := wrap(&api.ClientOutput{
			Message: message,
			IsDebug: false,
		})

		if err := stream.SendMsg(m); err != nil {
			logging.Log(s.Logger, "unable to write log to stream: %w", err)
		}
	}

	var debugTarget logging.Callback

	if options.CaptureDebugLog {
		debugTarget = func(f string, v ...interface{}) {
			send(fmt.Sprintf(f, v...), true)
		}
	}

	return &logging.CallbackLogger{
		LogTarget: func(f string, v ...interface{}) {
			send(fmt.Sprintf(f, v...), false)
		},
		DebugTarget: debugTarget,
	}
}
