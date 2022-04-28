package lambda

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/okhuz/bumper/models"
	"log"
	"strings"
)

func GetFunctionChan(lambda *models.Lambda) <-chan *models.Function {

	input := make(chan *models.Function, 100)

	go func() {
		for index := range lambda.Functions {
			if lambda.Functions[index].Bucket == "" {
				lambda.Functions[index].Bucket = lambda.Bucket
			}
			input <- &lambda.Functions[index]
		}
		close(input)
	}()

	return input
}

func GetUpdateFunctionChan(newFunction <-chan *models.Function) <-chan *models.Function {
	output := make(chan *models.Function, 100)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	lambdaClient := lambda.NewFromConfig(cfg)

	go func() {
		for f := range newFunction {
			updateFunction(s3Client, lambdaClient, f)
			output <- f
		}

		close(output)
	}()

	return output
}

func updateFunction(s3Client *s3.Client, lambdaClient *lambda.Client, newFunction *models.Function) {

	_, s3ObjectError := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &newFunction.Bucket, Key: &newFunction.Key})
	if s3ObjectError != nil {
		log.Printf("s3://%s/%s does not exist for function %s", newFunction.Bucket, newFunction.Key, newFunction.Name)
	} else {
		oldFunction, oldFunctionError := lambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{
			FunctionName: &newFunction.Name,
			Qualifier:    nil,
		})
		if oldFunctionError != nil {
			if strings.ContainsAny(oldFunctionError.Error(), "ResourceNotFoundException") {
				log.Printf("Function %s not found", newFunction.Name)
			}
		} else {
			log.Println("Updating function", *oldFunction.Configuration.FunctionName)
			_, updateFunctionError := lambdaClient.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{FunctionName: &newFunction.Name, S3Bucket: &newFunction.Bucket, S3Key: &newFunction.Key})
			if updateFunctionError != nil {
				log.Println(updateFunctionError)
			} else {
				log.Printf("%s has been updated", newFunction.Name)
			}
		}
	}
}
