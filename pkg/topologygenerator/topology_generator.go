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
