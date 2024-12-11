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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Networks-it-uc3m/l2sm-md/pkg/l2sminterface"
)

func main() {
	resource, output, input, _ := takeArguments()
	values, err := os.ReadFile(input)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		os.Exit(1)
	}

	crGenerator, err := l2sminterface.NewCRGenerator(resource)
	if err != nil {
		fmt.Printf("Specified resource not defined: %v\n", err)
		os.Exit(1)
	}

	err = crGenerator.AddValues(values)
	if err != nil {
		fmt.Printf("Failed to add new values to the resource: %v\n", err)
		os.Exit(1)
	}

	yamlFile, err := crGenerator.CreateResource()
	if err != nil {
		fmt.Printf("Failed to generate the resource file: %v\n", err)
		os.Exit(1)
	}
	if output != "" {
		os.WriteFile(output, yamlFile, 0644)

	} else {
		fmt.Println(yamlFile)
	}
}

func takeArguments() (l2sminterface.ResourceType, string, string, error) {

	output := flag.String("output", "", "directory where the ned settings are specified. Required")
	input := flag.String("input", "", "directory where the ned's neighbors  are specified. Required")
	resource := flag.String("resource", "Overlay", "name of the node the script is executed in. Required.")
	flag.Parse()

	switch {
	case *input == "":
		return l2sminterface.ResourceType("Overlay"), "", "", errors.New("node name is not defined")
	}

	resType := l2sminterface.ResourceType(*resource)
	return resType, *output, *input, nil
}
