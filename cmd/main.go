package cmd

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userService{})
	fmt.Println("gRPC server is running on port 50051")
	grpcServer.Serve(lis)
}
