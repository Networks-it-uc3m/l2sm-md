package operator

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetClusterCertificates(clusterConfig *rest.Config) (map[string][]byte, error) {

	clusterList := make(map[string][]byte)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return map[string][]byte{}, err
	}

	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{LabelSelector: "l2sm-cert"})
	if err != nil {
		return map[string][]byte{}, err
	}

	for _, secret := range secrets.Items {
		clusterList[secret.Labels["l2sm-cert"]] = secret.Data["cert-value"]
	}

	return clusterList, nil
}
func CreateCertificateSecrets(clusterConfig *rest.Config, namespace string, clusterName string, certificateData []byte) error {

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return fmt.Errorf("could not create new cluster client: %v", err)
	}

	// Define the secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-cert", clusterName),
			Labels: map[string]string{
				"l2sm-cert": clusterName,
			},
		},
		Data: map[string][]byte{
			"cert-value": certificateData,
		},
		Type: corev1.SecretTypeOpaque,
	}

	// Create the secret
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create secret: %v", err)
	}

	return nil
}
