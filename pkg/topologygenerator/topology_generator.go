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

package topologygenerator

import (
	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
)

func GenerateTopology(nodes []string) []*l2smmd.Link {
	numNodes := len(nodes)
	links := make([]*l2smmd.Link, 0, numNodes*(numNodes-1)/2)

	for i := 0; i < numNodes; i++ {
		for j := i + 1; j < numNodes; j++ {
			link := &l2smmd.Link{
				EndpointA: nodes[i],
				EndpointB: nodes[j],
			}
			links = append(links, link)
		}
	}

	return links
}
