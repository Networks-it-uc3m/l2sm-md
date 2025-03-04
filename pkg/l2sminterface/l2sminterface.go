package l2sminterface

import (
	"fmt"

	"github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"
)

func GetWorkloadPatchInstructions(networkName string) []*l2smmd.FieldPatch {

	patches := []*l2smmd.FieldPatch{
		{Path: "spec.template.metadata.labels.l2sm", Value: "true"},
		{Path: "spec.template.metadata.labels.l2sm/app", Value: "<workload-name>"},
		{Path: "spec.template.metadata.annotations.l2sm/networks", Value: networkName},
		{Path: "spec.template.spec.containers[0].env[0].name", Value: "DNS_NAME"},
		{Path: "spec.template.spec.containers[0].env[0].value", Value: fmt.Sprintf("<workload-name>.%s.%s.l2sm", networkName, "inter")},
	}
	return patches
}
