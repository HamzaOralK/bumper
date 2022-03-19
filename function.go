package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"strings"
)

func getFunctionChan(lambda *Lambda) <-chan *Function {

	input := make(chan *Function, 100)

	go func() {
		for num := range lambda.Functions {
			if lambda.Functions[num].Bucket == "" {
				lambda.Functions[num].Bucket = lambda.Bucket
			}
			input <- &lambda.Functions[num]
		}
		close(input)
	}()

	return input
}

func getUpdateFunctionChan(input <-chan *Function) <-chan *Function {
	output := make(chan *Function, 100)

	go func() {
		for f := range input {
			cfg, err := config.LoadDefaultConfig(context.TODO())
			if err != nil {
				panic(err)
			}

			s3Client := s3.NewFromConfig(cfg)
			_, s3ObjectError := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &f.Bucket, Key: &f.Key})
			if s3ObjectError != nil {
				log.Printf("s3://%s/%s does not exist for function %s", f.Bucket, f.Key, f.Name)
			} else {
				lambdaClient := lambda.NewFromConfig(cfg)
				lambdaInformation, lambdaInformationError := lambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{
					FunctionName: &f.Name,
					Qualifier:    nil,
				})
				if lambdaInformationError != nil {
					if strings.ContainsAny(lambdaInformationError.Error(), "ResourceNotFoundException") {
						log.Printf("Function %s not found", f.Name)
					}
				} else {
					log.Println("Updating function", *lambdaInformation.Configuration.FunctionName)
					_, updateFunctionError := lambdaClient.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{FunctionName: &f.Name, S3Bucket: &f.Bucket, S3Key: &f.Key})
					if updateFunctionError != nil {
						log.Println(updateFunctionError)
					}
				}
			}
			output <- f
		}

		close(output)
	}()

	return output
}
