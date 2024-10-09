package mdclient

import (
	"fmt"

	"context"
	"fmt"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/pb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type MDClient struct {
	ClusterConfigs []rest.Config
}

func (mdcli *MDClient) CreateNetwork(network *pb.L2Network) error {

	fmt.Printf("Creating network %s", network.GetName())

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	decUnstructured.Decode(yamlContent, nil, obj)

	for _, clusterConfig := range mdcli.ClusterConfigs {

		dynClient, err := dynamic.NewForConfig(&clusterConfig)
		if err != nil {
			return fmt.Errorf("Error contacting cluster %s: %v\n", clusterConfig.String(), err)
		}
		resource := l2smv1.GroupVersion.WithResource("l2networks")

		_, err = dynClient.Resource(resource).Namespace("default").Create(context.Background(), obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error creating resource: %v\n", err)
		}

		return nil
	}

	return nil
}

func (mdcli *MDClient) DeleteNetwork(network string) error {

	fmt.Printf("Deleting network %s", network)

	for _, clusterConfig := range mdcli.ClusterConfigs {

		dynClient, err := dynamic.NewForConfig(&clusterConfig)
		if err != nil {
			return fmt.Errorf("Error contacting cluster %s: %v\n", clusterConfig.String(), err)
		}

		_, err = dynClient.Resource(schema.GroupVersionResource{Group: "l2sm.l2sm.k8s.local", Version: "v1", Resource: "l2smnetworks"}).Namespace("default").Create(context.Background(), obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error creating resource: %v\n", err)
		}

		return nil
	}

	return nil

}
