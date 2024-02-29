package main

import (
	"log"

	"context"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func initK8sClient(cfg *config) *kubernetes.Clientset {
	var config *rest.Config
	var err error
	if cfg.KubeConfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfigPath)
		if err != nil {
			log.Fatalf("Error getting Kubernetes config: %v\n", err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error getting In cluster config: %v\n", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error getting clientset: %v\n", err)
	}
	return clientset
}

func initK8sSecretClient(client *kubernetes.Clientset, ns string) v1.SecretInterface {
	return client.CoreV1().Secrets(ns)

}

func createSecret(client v1.SecretInterface, token string) {
	secret := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openmetadata-bot-token",
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "openmetadata-initializer",
				"app.kubernetes.io/part-of":    "openmetadata",
			},
		},
		Type: coreV1.SecretType("Opaque"),
		Data: map[string][]byte{
			"token": []byte(token),
		},
	}
	result, err := client.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Unable to create secret: %s", err)
	}
	log.Printf("Secret created: %s", result.Name)
}
