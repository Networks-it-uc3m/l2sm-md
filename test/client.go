// client.go
//
// This is a simple client that interacts with the given gRPC server to create and delete a network.
// We assume that the server is running on localhost:50051.
//
// Note: Ensure that you have the appropriate generated gRPC code and protos accessible and correctly imported.
// The imports to "github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd" and gRPC libraries are illustrative.
//
// Usage:
//   go run client.go
//
// It will create a network with the specified constants and then delete it.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
)

// Constants for our example request
const (
	serverAddress      = "localhost:50051"
	networkName        = "pint-network"
	providerName       = "myProvider"
	providerDomain     = "myDomain.example.com"
	podCidr            = "10.244.0.0/16"
	networkType        = "vnet"
	clusterName        = "kind-worker-cluster-1"
	clusterBearerToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6InlxdlpuNmhhNW43X2FGS05ENWtjOFdyTklvV2VhNnQ2ODZiNGJCd3hJYUEifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzMzOTE3Mjk5LCJpYXQiOjE3MzM5MTM2OTksImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJsMnNtLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJsMnNtLWNvbnRyb2xsZXItbWFuYWdlciIsInVpZCI6IjAyZTI3YWI3LTIzYzAtNDFiYy1hOTA4LTQ1YTdkZDcwNDk1YiJ9fSwibmJmIjoxNzMzOTEzNjk5LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6bDJzbS1zeXN0ZW06bDJzbS1jb250cm9sbGVyLW1hbmFnZXIifQ.GWHQAOMP8UMUQB8PA9BOPW_UnGaXyGZcYdZHwctKMEMZr_3eofAkQuU_CaG-4JxbE95ldl4PwHeA-Ev8fbLDWS9sqSObxthRNqskDE2NTwYh6CAnacfzsZDE1IrQOZZS73oqbfLsbobb9NbF71pUsXz2gVrxnG08uxWHvL74D0hnQehZCHYcGKJZV6QwXA7pLN1g6yIY8JnUI5Xsyknylh-MD25gLPRe0-6HAk8okBMPRBLZ_fSWNLZ02GnqE6gSLwaKs1T9dn9NjGVUpBEtkZGx3RRI0cW077iw1uBImWxERVyeMkFi5tKN3pb5-C-1m0Ynk1G9RML04ohq-gAxMQ"
	clusterApiKey      = "https://0.0.0.0:41673"
)

func main() {
	// Set up a connection to the gRPC server.

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server at %s: %v", serverAddress, err)
	}
	defer conn.Close()

	// Create a client for our L2SMMultiDomainService
	client := l2smmd.NewL2SMMultiDomainServiceClient(conn)

	// Create a context with timeout to avoid blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clusters := []*l2smmd.Cluster{
		{
			Name: clusterName,
			RestConfig: &l2smmd.RestConfig{
				BearerToken: clusterBearerToken,
				ApiKey:      clusterApiKey,
			},
			Overlay: &l2smmd.Overlay{
				Provider: &l2smmd.Provider{
					Name:   providerName,
					Domain: providerDomain,
				},
				Nodes: []string{"node1", "node2"},
				Links: []*l2smmd.Link{
					{
						EndpointA: "node1",
						EndpointB: "node2",
					},
				},
			},
		},
	}
	// Create a network
	fmt.Println("Creating network...")
	createReq := &l2smmd.CreateNetworkRequest{
		Namespace: "l2sm-system",
		Network: &l2smmd.L2Network{
			Name: networkName,
			Provider: &l2smmd.Provider{
				Name:   providerName,
				Domain: providerDomain,
			},
			PodCidr:  podCidr,
			Type:     networkType,
			Clusters: clusters,
		},
	}

	createRes, err := client.CreateNetwork(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create network: %v", err)
	}
	fmt.Printf("CreateNetwork response: %s\n", createRes.GetMessage())

	// Here you might do additional operations, but for illustration we directly move to delete.

	// Delete the network
	fmt.Println("Deleting network...")
	deleteReq := &l2smmd.DeleteNetworkRequest{
		Network:   &l2smmd.L2Network{Name: networkName, Clusters: clusters},
		Namespace: "l2sm-system",
	}
	deleteRes, err := client.DeleteNetwork(ctx, deleteReq)
	if err != nil {
		log.Fatalf("Failed to delete network: %v", err)
	}
	fmt.Printf("DeleteNetwork response: %s\n", deleteRes.GetMessage())
}
