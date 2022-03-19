package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"strings"
)

func getDeploymentChan(kubernetes *Kubernetes) <-chan *Deployment {
	input := make(chan *Deployment, 100)

	go func() {
		for num := range kubernetes.Deployments {
			input <- &kubernetes.Deployments[num]
		}
		close(input)
	}()

	return input
}

func getUpdateDeploymentsChan(input <-chan *Deployment) <-chan *Deployment {
	output := make(chan *Deployment, 100)

	client := getClient()

	go func() {
		for f := range input {
			log.Println("Getting information for the deployment ", f.Name)
			d, deploymentGetError := client.AppsV1().Deployments(f.Namespace).Get(context.TODO(), f.Name, v1.GetOptions{})
			if deploymentGetError != nil {
				log.Println(deploymentGetError)
			} else {
				log.Printf("Updating deployment %s with version %s for Cluster %s\n", f.Name, f.Version, d.ClusterName)

				d.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", getImageRepository(d.Spec.Template.Spec.Containers[0].Image), f.Version)
				_, updateErr := client.AppsV1().Deployments(f.Namespace).Update(context.TODO(), d, v1.UpdateOptions{})
				if updateErr != nil {
					log.Println(updateErr)
				}
				output <- f
			}
		}
		close(output)
	}()

	return output
}

func getClient() *kubernetes.Clientset {
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
