package service

import (
	"SocialNetHL/internal/helper"
	"SocialNetHL/internal/store"
	"SocialNetHL/pkg/tokenservice"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	tokenservice.UnimplementedValidateTokenServer
	SessionStore store.SessionStore
}

func (s *Server) ValidateToken(ctx context.Context, request *tokenservice.ValidateTokenRequest) (*tokenservice.ValidateTokenResponse, error) {
	session, err := s.SessionStore.LoadSession(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	/*if len(session.UserId) == 0 {
		return nil, status.Error(codes.NotFound, "SESSION_NOT_FOUND")
	}*/

	resp := tokenservice.ValidateTokenResponse{
		Token:     session.Token,
		UserId:    session.UserId,
		Id:        session.Id,
		CreatedAt: session.CreatedAt,
	}
	return &resp, nil
}

func NewTokenServiceServer(store store.SessionStore) *Server {
	port := helper.GetEnvValue("RPC_SERVER_PORT", "50051")
	log.Printf("Starting gRPC server on port: %v", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	server := &Server{SessionStore: store}
	tokenservice.RegisterValidateTokenServer(s, server)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
	return server
}
