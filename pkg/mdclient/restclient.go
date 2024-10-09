package mdclient

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"context"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RestClient struct {
	ClusterConfigs []rest.Config
	Namespace      string
}

func (restcli *RestClient) CreateNetwork(network *l2smmd.L2Network) error {

	fmt.Printf("Creating network %s", network.GetName())

	l2network, _ := restcli.ConstructL2NetworkFromL2smmd(network)
	unstructuredL2network, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(l2network)
	unstructuredObj := &unstructured.Unstructured{Object: unstructuredL2network}

	for _, clusterConfig := range restcli.ClusterConfigs {

		dynClient, err := dynamic.NewForConfig(&clusterConfig)
		if err != nil {
			return fmt.Errorf("Error contacting cluster %s: %v\n", clusterConfig.String(), err)
		}
		resource := l2smv1.GroupVersion.WithResource("l2networks")

		_, err = dynClient.Resource(resource).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error creating resource: %v\n", err)
		}

		return nil
	}

	return nil
}

func (restcli *RestClient) DeleteNetwork(network string) error {

	fmt.Printf("Deleting network %s", network)

	for _, clusterConfig := range restcli.ClusterConfigs {

		dynClient, err := dynamic.NewForConfig(&clusterConfig)
		if err != nil {
			return fmt.Errorf("Error contacting cluster %s: %v\n", clusterConfig.String(), err)
		}
		resource := l2smv1.GroupVersion.WithResource("l2networks")

		err = dynClient.Resource(resource).Namespace(restcli.Namespace).Delete(context.Background(), network, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("Error creating resource: %v\n", err)
		}

		return nil
	}

	return nil

}

func (restcli *RestClient) ConstructL2NetworkFromL2smmd(network *l2smmd.L2Network) (*l2smv1.L2Network, error) {
	l2network := &l2smv1.L2Network{
		ObjectMeta: metav1.ObjectMeta{
			Name:      network.Name,
			Namespace: restcli.Namespace,
		},
		Spec: l2smv1.L2NetworkSpec{
			// Type:   l2smv1.NetworkType(network.Type),
			Config: &network.PodCidr,
			Provider: &l2smv1.ProviderSpec{
				Name:   network.Provider.Name,
				Domain: network.Provider.Domain,
			},
		},
	}
	return l2network, nil
}

func GetRestConfigs(absKubeconfigDirectory string) []rest.Config {
	kubeFiles, err := os.ReadDir(absKubeconfigDirectory)
	if err != nil {
		fmt.Printf("Couldn't get kube config files. err: %s", err)
	}
	clusterConfigs := readKubernetesConfigs(absKubeconfigDirectory, kubeFiles)
	return clusterConfigs
}

func readKubernetesConfigs(absKubeconfigDirectory string, configDirectories []fs.DirEntry) []rest.Config {

	var clusterConfigs []rest.Config

	for _, configEntry := range configDirectories {
		fmt.Println(filepath.Join(absKubeconfigDirectory, configEntry.Name()))
		kubeconfig := filepath.Join(absKubeconfigDirectory, configEntry.Name())
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clusterConfigs = append(clusterConfigs, *config)
	}

	return clusterConfigs

}
