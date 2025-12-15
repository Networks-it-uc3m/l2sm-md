package l2sminterface

import (
	"fmt"

	"github.com/Networks-it-uc3m/l2sc-es/api/v1/l2sces"
)

func GetWorkloadPatchInstructions(networkName string) []*l2sces.FieldPatch {

	patches := []*l2sces.FieldPatch{
		{Path: "spec.template.metadata.labels.l2sm", Value: "true"},
		{Path: "spec.template.metadata.labels.l2sm/app", Value: "<workload-name>"},
		{Path: "spec.template.metadata.annotations.l2sm/networks", Value: networkName},
		{Path: "spec.template.spec.containers[0].env[0].name", Value: "DNS_NAME"},
		{Path: "spec.template.spec.containers[0].env[0].value", Value: fmt.Sprintf("<workload-name>.%s.%s.l2sm", networkName, "inter")},
	}
	return patches
}
