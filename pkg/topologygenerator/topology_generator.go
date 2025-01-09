package topologygenerator

type Node struct {
	name      string
	ipAddress string
}
type Link struct {
	endpointA Node
	endpointB Node
}

func GenerateTopology(nodes ...Node) []Link {
	numNodes := len(nodes)
	links := make([]Link, numNodes*(numNodes-1)/2)
	nodeIndex := 0
	neighIndex := 1
	for i := 1; i <= len(links); i++ {
		links[i] = Link{endpointA: nodes[nodeIndex], endpointB: nodes[nodeIndex+neighIndex]}
		neighIndex++
		if numNodes == neighIndex {
			nodeIndex++
			neighIndex = nodeIndex + 1
		}
	}
	return links
}
