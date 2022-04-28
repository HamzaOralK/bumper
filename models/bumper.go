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

type Bumper struct {
	Lambda Lambda `yaml:"lambda"`
}
