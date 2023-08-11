package service

import (
	"SocialNetHL/internal/helper"
	"SocialNetHL/internal/store"
	"SocialNetHL/pkg/tokenservice"
	"context"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type Server struct {
	tokenservice.UnimplementedValidateTokenServer
	SessionStore store.SessionStore
	tracer       *tracesdk.TracerProvider
}

func (s *Server) ValidateToken(ctx context.Context, request *tokenservice.ValidateTokenRequest) (*tokenservice.ValidateTokenResponse, error) {
	// Extract TraceID from header
	md, _ := metadata.FromIncomingContext(ctx)
	traceIdString := md["x-trace-id"][0]
	// Convert string to byte array
	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		return nil, err
	}
	// Creating a span context with a predefined trace-id
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	// Embedding span config into the context
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	ctx, span := s.tracer.Tracer("social-app").Start(ctx, "ValidateToken")
	defer span.End()

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

func NewTokenServiceServer(store store.SessionStore, tracer *tracesdk.TracerProvider) *Server {
	port := helper.GetEnvValue("RPC_SERVER_PORT", "50051")
	log.Printf("Starting gRPC server on port: %v", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	server := &Server{SessionStore: store, tracer: tracer}
	tokenservice.RegisterValidateTokenServer(s, server)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
	return server
}
