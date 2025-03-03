package l2sminterface

import "github.com/Networks-it-uc3m/l2sm-md/api/v1/l2smmd"

func GetWorkloadPatchInstructions(networkName string) []*l2smmd.FieldPatch {

	patches := []*l2smmd.FieldPatch{
		{Path: "spec.template.metadata.labels.l2sm", Value: "true"},
		{Path: "spec.template.metadata.labels.l2sm/app", Value: "<workload-name>"},
		{Path: "spec.template.metadata.annotations.l2sm/networks", Value: networkName},
		// l2smmd.FieldPatch{Path: "a",Value: "b" },

	}
	return patches
}
