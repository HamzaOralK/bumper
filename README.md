# Description
Bumper is just a declarative version changer for lambdas. It just takes environment variables
for AWS authentication.

Lambdas require a s3 object addresses to update and kubernetes deployment require a new image tag. This tool does not
wait for updates just makes sure buckets, function and deployments exists.

Example yaml file:

```yaml
lambda:
  bucket: bucket
  functions:
    - name: function1
      key: function/lambda-v1.zip
```

Example run command: 
```shell
go run ./... ./test.yaml
```