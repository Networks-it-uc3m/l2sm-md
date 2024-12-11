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
	"fmt"
	"strings"

	l2smv1 "github.com/Networks-it-uc3m/L2S-M/api/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceType string

const (
	Overlay           ResourceType = "Overlay"
	NetworkEdgeDevice ResourceType = "NetworkEdgeDevice"
	L2Network         ResourceType = "L2Network"
)

type CRGenerator interface {
	CreateResource() ([]byte, error)
	AddValues([]byte) error
}

func NewCRGenerator(resource ResourceType) (CRGenerator, error) {

	switch resource {
	case Overlay:
		return &OverlayGenerator{}, nil
	case NetworkEdgeDevice:
	case L2Network:

	}
	return nil, fmt.Errorf("type %s not supported", resource)
}

// GetGVR returns the GroupVersionResource for the given resource type
func GetGVR(resource ResourceType) schema.GroupVersionResource {
	return l2smv1.GroupVersion.WithResource(getPluralResourceName(resource))
}

// getPluralResourceName converts a ResourceType to its plural form for GVR
func getPluralResourceName(resource ResourceType) string {
	switch resource {
	case Overlay:
		return "overlays"
	case NetworkEdgeDevice:
		return "networkedgedevices"
	case L2Network:
		return "l2networks"
	default:
		return strings.ToLower(string(resource)) + "s"
	}
}

// GetKind returns the Kind for the given resource type
func GetKind(resource ResourceType) string {
	return string(resource)
}
