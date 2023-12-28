package main

import (
	auth "auth/proto"
	models "auth/src/models"
	services "auth/src/services"
	utils "auth/src/utils"
	"context"
	"log"
	"net"

	"github.com/gofor-little/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server must implement the interfaces defined by the generated code from your proto file
type server struct {
	auth.UnimplementedAuthServiceServer
	userService *services.UserService // Add this line
}

func NewServer() *server {
	db, err := utils.InitDB() // Initialize the database
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	userService := services.NewUserService(db, "users") // Adjust "users" to your collection name

	return &server{
		userService: userService,
	}
}

func (s *server) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Convert the request to your user model
	user := &models.User{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	}

	// Use UserService to register the user
	if err := s.userService.RegisterUser(user); err != nil {
		// Convert the error into a gRPC status error
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

    return &auth.RegisterResponse{Id: user.ID}, nil
}

func (s *server) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	_,token, err := s.userService.LoginUser(in.Email, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	
    return &auth.LoginResponse{Token: token}, nil
}

func (s *server) GetSecretKey(ctx context.Context, in *auth.SecretKeyRequest) (*auth.SecretKeyResponse, error) {
	userID := in.UserId
	expDuration := in.GetExpDuration()
	isExp := in.GetIsExp()
	expTime := float64(expDuration)

	user := models.User{
		ID: userID,
	}

	token, err := utils.GenerateToken(user,isExp,&expTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	
    return &auth.SecretKeyResponse{Secretkey: token}, nil
}

func main() {
	// Load env file
	if err := env.Load(".env"); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv := NewServer()
	auth.RegisterAuthServiceServer(grpcServer, srv)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}