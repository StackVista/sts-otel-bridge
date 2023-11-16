package grpc

import "google.golang.org/grpc"

type Listener interface {
	Register(s *grpc.Server)
}
