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

func GetUpdateFunctionChan(input <-chan *models.Function) <-chan *models.Function {
	output := make(chan *models.Function, 100)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	lambdaClient := lambda.NewFromConfig(cfg)

	go func() {
		for f := range input {
			updateFunction(s3Client, lambdaClient, f)
			output <- f
		}

		close(output)
	}()

	return output
}

func updateFunction(s3Client *s3.Client, lambdaClient *lambda.Client, function *models.Function) {

	_, s3ObjectError := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &function.Bucket, Key: &function.Key})
	if s3ObjectError != nil {
		log.Printf("s3://%s/%s does not exist for function %s", function.Bucket, function.Key, function.Name)
	} else {
		lambdaInformation, lambdaInformationError := lambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{
			FunctionName: &function.Name,
			Qualifier:    nil,
		})
		if lambdaInformationError != nil {
			if strings.ContainsAny(lambdaInformationError.Error(), "ResourceNotFoundException") {
				log.Printf("Function %s not found", function.Name)
			}
		} else {
			log.Println("Updating function", *lambdaInformation.Configuration.FunctionName)
			_, updateFunctionError := lambdaClient.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{FunctionName: &function.Name, S3Bucket: &function.Bucket, S3Key: &function.Key})
			if updateFunctionError != nil {
				log.Println(updateFunctionError)
			}
		}
	}
}
