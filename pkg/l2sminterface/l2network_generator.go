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
	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
	"github.com/Networks-it-uc3m/l2sm-md/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConstructL2NetworkFromL2smmd(network *l2smmd.L2Network) (*l2smv1.L2Network, error) {

	l2network := &l2smv1.L2Network{
		TypeMeta: metav1.TypeMeta{
			Kind:       GetKind(L2Network), // Fix: Use the actual kind name, not the resource name
			APIVersion: l2smv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: network.Name,
		},
		Spec: l2smv1.L2NetworkSpec{
			Type:   l2smv1.NetworkType(utils.DefaultIfEmpty(network.Type, "vnet")),
			Config: &network.PodCidr,
			Provider: &l2smv1.ProviderSpec{
				Name:   network.Provider.Name,
				Domain: network.Provider.Domain,
			},
		},
	}
	return l2network, nil
}