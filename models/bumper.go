package models

type Function struct {
	Name   string `yaml:"name"`
	Bucket string `yaml:"bucket"`
	Key    string `yaml:"key"`
}

type Lambda struct {
	Bucket    string     `yaml:"bucket"`
	Functions []Function `yaml:"functions"`
}

type Deployment struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Tag       string `yaml:"tag"`
	Restart   bool   `yaml:"restart"`
}

type Kubernetes struct {
	Deployments []Deployment
}

type Bumper struct {
	Lambda     Lambda     `yaml:"lambda"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
}
