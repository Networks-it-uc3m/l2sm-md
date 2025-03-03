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

package main

import (
	"context"
	"fmt"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/l2sminterface"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/mdclient"
)

// server implements the L2SMMultiDomainServiceServer interface
type server struct {
	l2smmd.UnimplementedL2SMMultiDomainServiceServer
	mdclient.MDClient
}

// CreateNetwork calls a method from mdclient to create a network
func (s *server) CreateNetwork(ctx context.Context, req *l2smmd.CreateNetworkRequest) (*l2smmd.CreateNetworkResponse, error) {
	err := s.MDClient.CreateNetwork(req.GetNetwork(), req.GetNamespace())
	// Call the mdclient.CreateNetwork method (to be implemented later)
	if err != nil {
		return nil, fmt.Errorf("could not create network: %v", err)
	}
	return &l2smmd.CreateNetworkResponse{Message: "Network created successfully", Patches: l2sminterface.GetWorkloadPatchInstructions(req.GetNetwork().GetName())}, nil
}

// DeleteNetwork calls a method from mdclient to delete a network
func (s *server) DeleteNetwork(ctx context.Context, req *l2smmd.DeleteNetworkRequest) (*l2smmd.DeleteNetworkResponse, error) {
	// Call the mdclient.DeleteNetwork method (to be implemented later)
	err := s.MDClient.DeleteNetwork(req.GetNetwork(), req.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("could not delete network: %v", err)
	}
	return &l2smmd.DeleteNetworkResponse{Message: "Network deleted successfully"}, nil
}

func (s *server) CreateSlice(ctx context.Context, req *l2smmd.CreateSliceRequest) (*l2smmd.CreateSliceResponse, error) {
	err := s.MDClient.CreateSlice(req.GetSlice(), req.GetNamespace())

	if err != nil {
		return nil, fmt.Errorf("could now create slice: %v", err)
	}

	return &l2smmd.CreateSliceResponse{Message: "Slice created succesfully"}, nil
}

func (s *server) DeleteSlice(ctx context.Context, req *l2smmd.DeleteSliceRequest) (*l2smmd.DeleteSliceResponse, error) {
	// Call the mdclient.DeleteNetwork method (to be implemented later)
	err := s.MDClient.DeleteSlice(req.GetSlice(), req.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("could not delete network: %v", err)
	}
	return &l2smmd.DeleteSliceResponse{Message: "Slice deleted successfully"}, nil
}
