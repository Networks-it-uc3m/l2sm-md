package mdclient

import (
	"fmt"

	"l2sm.local/l2sm-md/pkg/pb"
)

func CreateNetwork(network *pb.L2Network) error {

	fmt.Printf("Received network %s", network.GetName())

	return nil
}

func DeleteNetwork(network string) error {

	fmt.Printf("Deleting network %s", network)

	return nil
}
