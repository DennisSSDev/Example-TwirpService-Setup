package server

import (
	"context"

	"github.com/dennisssdev/Example-TwirpService-Setup/rpc/example-service"
)

func (s *Server) SendText(ctx context.Context, request *example.SendTextRequest) (*example.SendTextResponse, error) {
	// Put the business logic in here
	return &example.SendTextResponse{
		Result: "OK",
	}, nil
}
