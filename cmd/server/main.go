// main.go

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/mdclient"
)

const (
	KUBECONFIGS_PATH = "/etc/l2sm/.kube/"
)

// server implements the L2SMMultiDomainServiceServer interface
type server struct {
	l2smmd.UnimplementedL2SMMultiDomainServiceServer
	mdclient.MDClient
}

// CreateNetwork calls a method from mdclient to create a network
func (s *server) CreateNetwork(ctx context.Context, req *l2smmd.CreateNetworkRequest) (*l2smmd.CreateNetworkResponse, error) {
	err := s.MDClient.CreateNetwork(req.Network)
	// Call the mdclient.CreateNetwork method (to be implemented later)
	if err != nil {
		return nil, err
	}
	return &l2smmd.CreateNetworkResponse{Message: "Network created successfully"}, nil
}

// DeleteNetwork calls a method from mdclient to delete a network
func (s *server) DeleteNetwork(ctx context.Context, req *l2smmd.DeleteNetworkRequest) (*l2smmd.DeleteNetworkResponse, error) {
	// Call the mdclient.DeleteNetwork method (to be implemented later)
	err := s.MDClient.DeleteNetwork(req.NetworkName)
	if err != nil {
		return nil, err
	}
	return &l2smmd.DeleteNetworkResponse{Message: "Network deleted successfully"}, nil
}

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	restConfigs, err := mdclient.GetRestConfigs(KUBECONFIGS_PATH)
	if err != nil {
		log.Fatalf("Failed to get rest configs: %v", err)
	}

	restcli, err := mdclient.NewClient(mdclient.RestType, restConfigs)

	if err != nil {
		log.Fatalf("Failed to create multi domain client: %v", err)
	}

	// Register the server with the gRPC server
	l2smmd.RegisterL2SMMultiDomainServiceServer(grpcServer, &server{MDClient: restcli})

	log.Printf("Server listening at %v", lis.Addr())

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
