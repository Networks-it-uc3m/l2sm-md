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

package l2sminterface

import (
	"testing"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestApplyCIDRs_Enhanced runs multiple sub-tests to verify ApplyCIDRs functionality,
// including handling of too many clusters, a case with 3 clusters, zero clusters, and a /0 network.
func TestApplyCIDRs_Enhanced(t *testing.T) {
	// Prepare a base L2Network (with minimal required fields).
	baseL2Network := &l2smv1.L2Network{
		TypeMeta: metav1.TypeMeta{
			Kind:       GetKind(L2Network),
			APIVersion: l2smv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: l2smv1.L2NetworkSpec{
			Type: "vnet",
		},
	}

	tests := []struct {
		name           string
		networkCIDR    string
		numberClusters int
		expectedCIDRs  []string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Valid with 5 clusters",
			networkCIDR:    "192.168.0.0/16",
			numberClusters: 5,
			// For 5 clusters, bitsNeeded = len("101") = 3, so newPrefix = 16+3 = 19.
			// Block size = 2^(32-19) = 8192 addresses (32 in the third octet).
			expectedCIDRs: []string{
				"192.168.32.0/19",
				"192.168.64.0/19",
				"192.168.96.0/19",
				"192.168.128.0/19",
				"192.168.160.0/19",
			},
			expectError: false,
		},
		{
			name:           "Valid with 3 clusters",
			networkCIDR:    "192.168.0.0/16",
			numberClusters: 3,
			// For 3 clusters, bitsNeeded = len("11") = 2, so newPrefix = 16+2 = 18.
			// Block size = 2^(32-18)=16384 addresses.
			// Subnet addresses:
			//   i=1: 192.168.0.0 + 16384 = 192.168.64.0/18,
			//   i=2: 192.168.0.0 + 32768 = 192.168.128.0/18,
			//   i=3: 192.168.0.0 + 49152 = 192.168.192.0/18.
			expectedCIDRs: []string{
				"192.168.64.0/18",
				"192.168.128.0/18",
				"192.168.192.0/18",
			},
			expectError: false,
		},
		{
			name:           "Too many clusters causing prefix > 32",
			networkCIDR:    "192.168.0.0/16",
			numberClusters: 100000, // 100000 in binary is 17 digits, so newPrefix = 16+17 = 33.
			expectedCIDRs:  nil,
			expectError:    true,
			expectedErrMsg: "new prefix length 33 exceeds 32",
		},
		{
			name:           "Zero clusters returns empty list",
			networkCIDR:    "192.168.0.0/16",
			numberClusters: 0,
			expectedCIDRs:  []string{},
			expectError:    false,
		},
		{
			name:           "Edge case with /0 network and 3 clusters",
			networkCIDR:    "0.0.0.0/0",
			numberClusters: 3,
			// For /0, ones = 0, bitsNeeded for 3 = 2 so newPrefix = 2.
			// Block size = 2^(32-2)=2^30 = 1073741824.
			// Base IP is 0.0.0.0, so:
			//   i=1: 0 + 1073741824 -> 64.0.0.0/2,
			//   i=2: 2*1073741824 -> 128.0.0.0/2,
			//   i=3: 3*1073741824 -> 192.0.0.0/2.
			expectedCIDRs: []string{
				"64.0.0.0/2",
				"128.0.0.0/2",
				"192.0.0.0/2",
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ApplyCIDRs(tc.networkCIDR, *baseL2Network, tc.numberClusters)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tc.expectedErrMsg {
					t.Fatalf("expected error message %q, got %q", tc.expectedErrMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.Items) != len(tc.expectedCIDRs) {
				t.Fatalf("expected %d subnets, got %d", len(tc.expectedCIDRs), len(result.Items))
			}

			for i, netw := range result.Items {
				// Here we compare the CIDR field; adjust this if your struct stores it elsewhere.
				if netw.Spec.NetworkCIDR != tc.expectedCIDRs[i] {
					t.Errorf("expected CIDR %s, got %s", tc.expectedCIDRs[i], netw.Spec.NetworkCIDR)
				}
			}
		})
	}
}
