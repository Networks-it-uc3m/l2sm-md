package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Networks-it-uc3m/l2sm-md/pkg/operator"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: program [options] <certificate_file_path>")
		os.Exit(1)
	}

	// Parse command-line flags
	var kubeconfig *string
	var namespace *string
	var clusterName *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespace = flag.String("namespace", "default", "Kubernetes namespace")
	clusterName = flag.String("clustername", "test", "Cluster name")
	flag.Parse()

	// Get the certificate file path from the last argument
	certificateFilePath := flag.Arg(0)
	if certificateFilePath == "" {
		fmt.Println("Certificate file path is required as the last argument.")
		os.Exit(1)
	}

	// Read the certificate data from the file
	certificate, err := ioutil.ReadFile(certificateFilePath)
	if err != nil {
		fmt.Printf("Failed to read certificate file: %v\n", err)
		os.Exit(1)
	}

	// Build Kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create the secret
	err = operator.CreateCertificateSecrets(config, *namespace, *clusterName, certificate)
	if err != nil {
		panic(err)
	}

	fmt.Println("Secret created successfully.")
}
