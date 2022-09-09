package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CustomResourceEntryArray []CustomResourceEntry `yaml:"crds"`
}

type CustomResourceEntry struct {
	Api        string `yaml:"api"`
	ApiVersion string `yaml:"apiversion"`
	Namespace  string `yaml:"namespace"`
	Resource   string `yaml:"resource"`
	Name       string `yaml:"name"`
}

func ValidateConfigFilePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return errors.New(fmt.Sprintf("Path: %s is a directory. Please pass a direct file path.", path))
	}
	return nil
}

func NewConfig(path string) (*Config, error) {
	conf := new(Config)
	//open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() //keep file open until the function returns.
	yamlContent := yaml.NewDecoder(file)
	decodeErr := yamlContent.Decode(&conf)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return conf, nil
}
