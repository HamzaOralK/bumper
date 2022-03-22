package kubernetes

import (
	"context"
	"fmt"
	"github.com/okhuz/bumper/models"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"strings"
)

func GetDeploymentChan(kubernetes *models.Kubernetes) <-chan *models.Deployment {
	input := make(chan *models.Deployment, 100)

	go func() {
		for num := range kubernetes.Deployments {
			input <- &kubernetes.Deployments[num]
		}
		close(input)
	}()

	return input
}

func GetUpdateDeploymentsChan(input <-chan *models.Deployment) <-chan *models.Deployment {
	output := make(chan *models.Deployment, 100)
	client := getKubernetesClientSet()

	go func() {
		for f := range input {
			_ = updateDeployment(client, f)
			output <- f
		}
		close(output)
	}()

	return output
}

func updateDeployment(client *kubernetes.Clientset, deployment *models.Deployment) error {
	log.Println("Getting information for the deployment ", deployment.Name)
	deployed, deploymentGetError := client.AppsV1().Deployments(deployment.Namespace).Get(context.TODO(), deployment.Name, v1.GetOptions{})
	if deploymentGetError != nil {
		log.Println(deploymentGetError)
		return deploymentGetError
	} else {
		log.Printf("Updating deployment %s with version %s for Cluster %s\n", deployment.Name, deployment.Tag, deployed.ClusterName)
		deployed.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", getImageRepository(deployed.Spec.Template.Spec.Containers[0].Image), deployment.Tag)
		_, updateErr := client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployed, v1.UpdateOptions{})
		if updateErr != nil {
			log.Println(updateErr)
		}
		return nil
	}
}

func getKubernetesClientSet() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet
}

func getImageRepository(image string) string {
	return strings.Split(image, ":")[0]
}
