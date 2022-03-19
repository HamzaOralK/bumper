# Description
Bumper is just a declarative version changer for lambdas and Kubernetes cluster. It just takes environment variables and
kubeconfig to connect Kubernetes cluster.

Lambdas require a s3 object addresses to update and kubernetes deployment require a new image tag. This tool does not
wait for updates just makes sure everything exists.

Example yaml file:

```yaml
lambda:
  bucket: bucket
  functions:
    - name: function1
      key: function/lambda-v1.zip

kubernetes:
  deployments:
    - name: service1
      namespace: default
      version: v1
    - name: service2
      namespace: tenera
      version: v2
    - name: service3
      namespace: tenera
      version: v3
```

Example run command: 
```shell
go run ./... ./test.yaml
```