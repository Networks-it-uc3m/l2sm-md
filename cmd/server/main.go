// main.go

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Networks-it-uc3m/l2sm-md/pkg/mdclient"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/pb"
)

// server implements the L2SMMultiDomainServiceServer interface
type server struct {
	pb.UnimplementedL2SMMultiDomainServiceServer
	mdclient.MDClient
}

// CreateNetwork calls a method from mdclient to create a network
func (s *server) CreateNetwork(ctx context.Context, req *pb.CreateNetworkRequest) (*pb.CreateNetworkResponse, error) {
	err := s.MDClient.CreateNetwork(req.Network)
	// Call the mdclient.CreateNetwork method (to be implemented later)
	if err != nil {
		return nil, err
	}
	return &pb.CreateNetworkResponse{Message: "Network created successfully"}, nil
}

// DeleteNetwork calls a method from mdclient to delete a network
func (s *server) DeleteNetwork(ctx context.Context, req *pb.DeleteNetworkRequest) (*pb.DeleteNetworkResponse, error) {
	// Call the mdclient.DeleteNetwork method (to be implemented later)
	err := s.MDClient.DeleteNetwork(req.NetworkName)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteNetworkResponse{Message: "Network deleted successfully"}, nil
}

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the server with the gRPC server
	pb.RegisterL2SMMultiDomainServiceServer(grpcServer, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
