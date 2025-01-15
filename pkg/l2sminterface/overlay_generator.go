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
	"encoding/json"
	"fmt"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/topologygenerator"
	"gopkg.in/yaml.v2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type OverlayGenerator struct {
	Values *l2smv1.TopologySpec
}

func constructOverlayFromTopology(overlay *l2smv1.TopologySpec) (*l2smv1.Overlay, error) {

	l2smOverlay := &l2smv1.Overlay{
		TypeMeta: metav1.TypeMeta{
			APIVersion: l2smv1.GroupVersion.Identifier(),
			Kind:       "Overlay",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "overlay-sample",
		},
		Spec: l2smv1.OverlaySpec{
			NetworkController: defaultNetworkController(),

			SwitchTemplate: defaultSwitchTemplate(),
			Topology: &l2smv1.TopologySpec{
				Nodes: overlay.Nodes,
				Links: overlay.Links,
			},
		},
	}
	return l2smOverlay, nil
}
func (overlayGenerator *OverlayGenerator) CreateResource() ([]byte, error) {
	l2smOverlay, err := constructOverlayFromTopology(overlayGenerator.Values)
	if err != nil {
		return nil, fmt.Errorf("could not construct overlay, given the input values. Error: %v", err)
	}

	// Convert the structured object to an unstructured one
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(l2smOverlay)
	if err != nil {
		return nil, fmt.Errorf("could not convert to unstructured: %v", err)
	}

	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

	// Remove creationTimestamp field from metadata
	if metadata, ok := unstructuredObj.Object["metadata"].(map[string]interface{}); ok {
		delete(metadata, "creationTimestamp")
	}

	// Remove creationTimestamp field from spec.switchTemplate.metadata
	if spec, ok := unstructuredObj.Object["spec"].(map[string]interface{}); ok {
		if switchTemplate, ok := spec["switchTemplate"].(map[string]interface{}); ok {
			if switchTemplateMetadata, ok := switchTemplate["metadata"].(map[string]interface{}); ok {
				delete(switchTemplateMetadata, "creationTimestamp")
			}
		}
	}

	// Marshal the unstructured object to YAML
	yamlData, err := yaml.Marshal(unstructuredObj.Object)
	if err != nil {
		return nil, fmt.Errorf("could not marshal to YAML: %v", err)
	}

	return yamlData, nil
}

// func (overlayGenerator *OverlayGenerator) AddValues(byteValues []byte) error {

// 	values := l2smmd.Overlay{}
// 	err := yaml.Unmarshal(byteValues, &values)
// 	if err != nil {
// 		return fmt.Errorf("could not unmarshal input values. err: %v", err)
// 	}
// 	overlayGenerator.Values = values
// 	return nil
// }

func (overlayGenerator *OverlayGenerator) AddValues(byteValues []byte) error {
	// Create an instance of l2smmd.Overlay to hold the unmarshaled values
	values := l2smv1.TopologySpec{}

	// Use yaml.Unmarshal to populate the values, passing its pointer
	err := json.Unmarshal(byteValues, &values)
	if err != nil {
		return fmt.Errorf("could not unmarshal input values. err: %v", err)
	}

	// Assign the unmarshaled values to the overlayGenerator.Values field
	overlayGenerator.Values = &values
	return nil
}

func ConstructOverlayFromL2smmd(overlay *l2smmd.Overlay) *l2smv1.Overlay {
	fmt.Println("OVERLAYY")
	fmt.Println(overlay.GetProvider().GetName())
	fmt.Println("safe?")

	links := make([]l2smv1.Link, len(overlay.Links))

	if len(overlay.Links) == 0 && len(overlay.Nodes) > 1 {
		overlay.Links = topologygenerator.GenerateTopology(overlay.GetNodes())
	}

	for _, link := range overlay.Links {
		l2Link := l2smv1.Link{EndpointA: link.EndpointA, EndpointB: link.EndpointB}
		links = append(links, l2Link)
	}

	l2overlay := &l2smv1.Overlay{
		TypeMeta: metav1.TypeMeta{
			Kind:       GetKind(Overlay), // Fix: Use the actual kind name, not the resource name
			APIVersion: l2smv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "overlay-sample",
		},
		Spec: l2smv1.OverlaySpec{
			NetworkController: defaultNetworkController(),
			SwitchTemplate:    defaultSwitchTemplate(),
			Topology: &l2smv1.TopologySpec{
				Nodes: overlay.Nodes,
				Links: links,
			},
		},
	}
	return l2overlay
}

func defaultSwitchTemplate() *l2smv1.SwitchTemplateSpec {
	return &l2smv1.SwitchTemplateSpec{
		Spec: l2smv1.SwitchPodSpec{
			Containers: []corev1.Container{
				{
					Name:  "l2sm-switch",
					Image: "alexdecb/l2sm-switch:2.7",
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
					ImagePullPolicy: corev1.PullAlways,
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

func defaultNetworkController() *l2smv1.NetworkControllerSpec {
	return &l2smv1.NetworkControllerSpec{
		Name:   "l2sm-sdn",
		Domain: "l2sm-controller-service.l2sm-system.svc",
	}
}
