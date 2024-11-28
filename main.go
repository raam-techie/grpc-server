package main

import (
	"context"
	"log"
	"net"
	"os"

	mysvccore "grpc-server/core"
	mysvcgrpc "grpc-server/grpc-server/grpc/pb"
	"grpc-server/mysvc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting the server...")

	// configure our core service
	userService := mysvccore.NewService()

	// configure our gRPC service controller
	userServiceController := NewUserServiceController(userService)

	// start a gRPC server
	server := grpc.NewServer()
	log.Println("gRPC server created")

	mysvcgrpc.RegisterUserServiceServer(server, userServiceController)

	reflection.Register(server)

	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "localhost:8080" // Default address
	}

	log.Printf("Binding to address: %s\n", addr)
	con, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v\n", addr, err)
	}

	log.Println("Starting server...")
	err = server.Serve(con)

	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}

	log.Println("Application up and running on port 8080")
}

// userServiceController implements the gRPC UserServiceServer interface.
type userServiceController struct {
	mysvcgrpc.UnimplementedUserServiceServer // Embed the auto-generated unimplemented server
	userService                              mysvc.Service
}

// NewUserServiceController instantiates a new UserServiceServer.
func NewUserServiceController(userService mysvc.Service) mysvcgrpc.UserServiceServer {
	log.Println("Creating UserServiceController...")
	return &userServiceController{
		userService: userService,
	}
}

// GetUsers calls the core service's GetUsers method and maps the result to a grpc service response.
func (ctlr *userServiceController) GetUsers(ctx context.Context, req *mysvcgrpc.GetUsersRequest) (resp *mysvcgrpc.GetUsersResponse, err error) {
	log.Println("Received GetUsers request...")
	resultMap, err := ctlr.userService.GetUsers(req.GetIds())
	if err != nil {
		log.Printf("Error in GetUsers: %v\n", err)
		return
	}
	resp = &mysvcgrpc.GetUsersResponse{}
	for _, u := range resultMap {
		resp.Users = append(resp.Users, marshalUser(&u))
	}
	return
}

// marshalUser marshals a business object User into a gRPC layer User.
func marshalUser(u *mysvc.User) *mysvcgrpc.User {
	return &mysvcgrpc.User{Id: u.ID, Name: u.Name}
}
