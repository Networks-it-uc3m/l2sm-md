package mdclient

import (
	"errors"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"k8s.io/client-go/rest"
)

type ClientType string

const (
	RestType ClientType = "rest"
)

type MDClient interface {
	CreateNetwork(network *l2smmd.L2Network) error
	DeleteNetwork(network string) error
}

func NewClient(clientType ClientType, config ...interface{}) (MDClient, error) {

	switch clientType {
	case RestType:
		clusterConfigs := []rest.Config{}
		// Convert each element in the config slice to rest.Config
		for _, cfg := range config {
			// Assert that cfg is of type rest.Config
			if c, ok := cfg.(rest.Config); ok {
				clusterConfigs = append(clusterConfigs, c)
			} else {
				return nil, errors.New("invalid config type, expected rest.Config")
			}
		}
		client := &RestClient{ClusterConfigs: clusterConfigs}
		return client, nil
	default:
		return nil, errors.New("unsupported client type")
	}
}
