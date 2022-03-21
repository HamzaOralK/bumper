package main

import (
	"github.com/okhuz/bumper/kubernetes"
	"github.com/okhuz/bumper/lambda"
	"github.com/okhuz/bumper/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

func main() {
	arguments := os.Args[1:]

	file, err := os.ReadFile(arguments[0])
	if err != nil {
		log.Fatal(err.Error())
	}
	bumper := models.Bumper{}

	err = yaml.Unmarshal(file, &bumper)
	if err != nil {
		log.Fatal(err.Error())
	}

	functionChan := lambda.GetFunctionChan(&bumper.Lambda)
	functionOperationChan := lambda.GetUpdateFunctionChan(functionChan)

	deploymentChan := kubernetes.GetDeploymentChan(&bumper.Kubernetes)
	deploymentOperationChan := kubernetes.GetUpdateDeploymentsChan(deploymentChan)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range functionOperationChan {
			// fmt.Println(f)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range deploymentOperationChan {
			// fmt.Println(d)
		}
	}()

	wg.Wait()

}
