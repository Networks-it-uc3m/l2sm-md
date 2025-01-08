// Copyright 2024 Universidad Carlos III de Madrid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the gRPC server.

	// Define a boolean flag for test slice creation
	testNetworkCreate := flag.Bool("test-network-create", false, "Simulate the creation and deployment of network resources")
	testSliceCreate := flag.Bool("test-slice-create", false, "Simulate the creation and deployment of slice resources")
	testSliceDelete := flag.Bool("test-slice-delete", false, "Simulate the deletion of slice resources")
	testNetworkDelete := flag.Bool("test-network-delete", false, "Simulate the deletion of network resources")
	configPath := flag.String("config", "./config.yaml", "path to config file")

	// Parse the command-line flags
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err != nil {
		log.Fatalf("Failed to connect to server at %s: %v", cfg.ServerAddress, err)
	}
	defer conn.Close()

	// Create a client for our L2SMMultiDomainService
	client := l2smmd.NewL2SMMultiDomainServiceClient(conn)

	clusters := make([]*l2smmd.Cluster, 0, len(cfg.Clusters))
	for _, c := range cfg.Clusters {
		clusters = append(clusters, &l2smmd.Cluster{
			Name: c.Name,
			RestConfig: &l2smmd.RestConfig{
				BearerToken: c.BearerToken,
				ApiKey:      c.ApiKey,
			},
			Overlay: &l2smmd.Overlay{
				Nodes: c.Nodes,
			},
		})
	}

	if *testSliceCreate {

		fmt.Println("Creating network...")
		createSliceReq := &l2smmd.CreateSliceRequest{
			Namespace: "l2sm-system",
			Slice: &l2smmd.Slice{
				Provider: &l2smmd.Provider{
					Name:   cfg.Provider.Name,
					Domain: cfg.Provider.Domain,
				},
				Clusters: clusters,
			},
		}
		_, err := client.CreateSlice(ctx, createSliceReq)
		if err != nil {
			log.Fatalf("Failed to create slice: %v", err)
		}
		fmt.Printf("CreateNetwork response: %s\n", "Network 'ping-network' created successfully.")
	}

	if *testSliceDelete {
	}
	if *testNetworkDelete {
	}
	if *testNetworkCreate {
		fmt.Println("Creating network...")

		// Create a context with timeout to avoid blocking indefinitely
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Create a network
		fmt.Println("Creating network...")
		createReq := &l2smmd.CreateNetworkRequest{
			Namespace: "l2sm-system",
			Network: &l2smmd.L2Network{
				Name: cfg.NetworkName,
				Provider: &l2smmd.Provider{
					Name:   cfg.Provider.Name,
					Domain: cfg.Provider.Domain,
				},
				Type:     cfg.NetworkType,
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
			Network:   &l2smmd.L2Network{Name: cfg.NetworkName, Clusters: clusters},
			Namespace: "l2sm-system",
		}
		deleteRes, err := client.DeleteNetwork(ctx, deleteReq)
		if err != nil {
			log.Fatalf("Failed to delete network: %v", err)
		}
		fmt.Printf("DeleteNetwork response: %s\n", deleteRes.GetMessage())

	}
}
