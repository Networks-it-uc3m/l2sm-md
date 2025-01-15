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
	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NEDValues struct {
	NodeConfig NodeConfig
	Neighbors  []Neighbor
}

type SDNController struct {
	Name   string
	Domain string
}

type NodeConfig struct {
	NodeName  string
	IPAddress string
}

type Neighbor struct {
	Node   string
	Domain string
}

type NEDGenerator struct {
	SliceName string
	Provider  SDNController
}

func NewNEDGenerator(sliceName string, providerDomain string) *NEDGenerator {
	return &NEDGenerator{
		SliceName: sliceName,
		Provider: SDNController{
			Name:   sliceName + "-controller",
			Domain: providerDomain,
		}}
}
func (nedGenerator *NEDGenerator) ConstructNED(nedValues NEDValues) *l2smv1.NetworkEdgeDevice {

	neighbors := make([]l2smv1.NeighborSpec, len(nedValues.Neighbors))
	for i := range neighbors {
		neighbors[i].Node = nedValues.Neighbors[i].Node
		neighbors[i].Domain = nedValues.Neighbors[i].Domain
	}
	ned := &l2smv1.NetworkEdgeDevice{
		TypeMeta: metav1.TypeMeta{
			Kind:       GetKind(NetworkEdgeDevice),
			APIVersion: l2smv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: nedGenerator.SliceName + "-ned",
		},
		Spec: l2smv1.NetworkEdgeDeviceSpec{
			NetworkController: &l2smv1.NetworkControllerSpec{
				Name:   nedGenerator.Provider.Name,
				Domain: nedGenerator.Provider.Domain,
			},
			NodeConfig: &l2smv1.NodeConfigSpec{
				NodeName:  nedValues.NodeConfig.NodeName,
				IPAddress: nedValues.NodeConfig.IPAddress,
			},
			Neighbors:      neighbors,
			SwitchTemplate: defaultNEDTemplate(),
		},
	}
	return ned

}

func defaultNEDTemplate() *l2smv1.SwitchTemplateSpec {
	return &l2smv1.SwitchTemplateSpec{
		Spec: l2smv1.SwitchPodSpec{
			HostNetwork: true,
			Containers: []corev1.Container{
				{
					Name:  "l2sm-ned",
					Image: "alexdecb/l2sm-switch:2.7.1",
					Env: []corev1.EnvVar{
						{
							Name: "NODENAME",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "spec.nodeName",
								},
							},
						},
					},
					SecurityContext: &corev1.SecurityContext{
						Capabilities: &corev1.Capabilities{
							Add: []corev1.Capability{"NET_ADMIN"},
						},
					},
				},
			},
		},
	}
}
