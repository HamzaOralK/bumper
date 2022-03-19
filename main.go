package main

import (
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
	bumper := Bumper{}

	err = yaml.Unmarshal(file, &bumper)
	if err != nil {
		log.Fatal(err.Error())
	}

	functionChan := getFunctionChan(&bumper.Lambda)
	functionOperationChan := getUpdateFunctionChan(functionChan)

	deploymentChan := getDeploymentChan(&bumper.Kubernetes)
	deploymentOperationChan := getUpdateDeploymentsChan(deploymentChan)

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
