package main

import (
	"context"
	"errors"
	"log"
	"net"

	firebase "firebase.google.com/go"
	pb "github.com/ahmadexe/prism-backend/grpc/auth/generated"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type AuthGrpcServer struct {
	pb.AuthServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, &AuthGrpcServer{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *AuthGrpcServer) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	opt := option.WithCredentialsFile("../../env_var/app_keys.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return &pb.AuthorizeResponse{
			IsAuthorized: false,
		}, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return &pb.AuthorizeResponse{
			IsAuthorized: false,
		}, err
	}
	
	_, err = client.VerifyIDToken(ctx, req.Token)
	if err != nil {
		en := errors.New("invalid token")

		return &pb.AuthorizeResponse{
			IsAuthorized: false,
		}, en
	}

	return &pb.AuthorizeResponse{
		IsAuthorized: true,
	}, nil
}