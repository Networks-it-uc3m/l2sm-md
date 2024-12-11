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

package mdclient

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"context"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/operator"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RestClient struct {
	ManagerClusterConfig rest.Config
}

func (restcli *RestClient) CreateNetwork(network *l2smmd.L2Network, namespace string) error {

	fmt.Printf("Creating network %s", network.GetName())
	namespace = utils.DefaultIfEmpty(namespace, "default")

	l2network, err := l2sminterface.ConstructL2NetworkFromL2smmd(network)
	if err != nil {
		return fmt.Errorf("failed to construct l2network: %v", err)
	}
	unstructuredL2network, err := runtime.DefaultUnstructuredConverter.ToUnstructured(l2network)
	if err != nil {
		return fmt.Errorf("failed to assign unstructured l2network: %v", err)
	}
	unstructuredObj := &unstructured.Unstructured{Object: unstructuredL2network}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("could not create cluster config for control plane cluster: %v", err)

	}
	clusterCrts, err := operator.GetClusterCertificates(config)

	if err != nil {
		return fmt.Errorf("could not get cluster certificates error: %v", err)
	}

	for _, cluster := range network.Clusters {

		clusterConfig := &rest.Config{Host: cluster.RestConfig.ApiKey, BearerToken: cluster.RestConfig.BearerToken,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: false, // Set to true if self-signed certs are acceptable
				CAData:   clusterCrts[cluster.Name],
			},
		}
		dynClient, err := dynamic.NewForConfig(clusterConfig)
		if err != nil {
			return fmt.Errorf("error contacting cluster %s: %v", clusterConfig.String(), err)
		}

		resource := l2sminterface.GetGVR(l2sminterface.L2Network)

		_, err = dynClient.Resource(resource).Namespace(namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating resource: %v", err)
		}

	}

	return nil
}

func (restcli *RestClient) DeleteNetwork(network *l2smmd.L2Network, namespace string) error {

	clusterCrts, err := operator.GetClusterCertificates(&restcli.ManagerClusterConfig)

	if err != nil {
		return fmt.Errorf("could not get cluster certificates error: %v", err)
	}

	fmt.Printf("Deleting network %s", network.Name)
	namespace = utils.DefaultIfEmpty(namespace, "default")
	for _, cluster := range network.Clusters {
		clusterConfig := &rest.Config{Host: cluster.RestConfig.ApiKey, BearerToken: cluster.RestConfig.BearerToken,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: false, // Set to true if self-signed certs are acceptable
				CAData:   clusterCrts[cluster.Name],
			},
		}
		dynClient, err := dynamic.NewForConfig(clusterConfig)
		if err != nil {
			return fmt.Errorf("error contacting cluster %s: %v", clusterConfig.String(), err)
		}
		resource := l2sminterface.GetGVR(l2sminterface.L2Network)

		err = dynClient.Resource(resource).Namespace(namespace).Delete(context.Background(), network.Name, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("error deleting resource: %v", err)
		}

	}

	return nil

}

func GetRestConfigs(absKubeconfigDirectory string) ([]rest.Config, error) {
	kubeFiles, err := os.ReadDir(absKubeconfigDirectory)
	if err != nil {
		return []rest.Config{}, fmt.Errorf("couldn't get kube config files in %s: %v", absKubeconfigDirectory, err)

	}
	clusterConfigs, err := readKubernetesConfigs(absKubeconfigDirectory, kubeFiles)
	if err != nil {
		return []rest.Config{}, fmt.Errorf("failed to read configs in %s: %v", absKubeconfigDirectory, err)
	}
	return clusterConfigs, nil
}

func readKubernetesConfigs(absKubeconfigDirectory string, configDirectories []fs.DirEntry) ([]rest.Config, error) {

	var clusterConfigs []rest.Config

	for _, configEntry := range configDirectories {
		fmt.Println(filepath.Join(absKubeconfigDirectory, configEntry.Name()))
		kubeconfig := filepath.Join(absKubeconfigDirectory, configEntry.Name())
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return []rest.Config{}, fmt.Errorf("failed to build config %s from flags: %v", configEntry.Name(), err)
		}
		clusterConfigs = append(clusterConfigs, *config)
	}

	return clusterConfigs, nil

}
