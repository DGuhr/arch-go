package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type DependenciesRule struct {
	Package             string   `yaml:"package"`
	ShouldOnlyDependsOn []string `yaml:"shouldOnlyDependsOn"`
	ShouldNotDependsOn  []string `yaml:"shouldNotDependsOn"`
}

type ContentsRule struct {
	Package                     string `yaml:"package"`
	ShouldOnlyContainInterfaces bool   `yaml:"shouldOnlyContainInterfaces"`
	ShouldOnlyContainTypes      bool   `yaml:"shouldOnlyContainTypes"`
	ShouldOnlyContainFunctions  bool   `yaml:"shouldOnlyContainFunctions"`
	ShouldOnlyContainMethods    bool   `yaml:"shouldOnlyContainMethods"`
	ShouldNotContainInterfaces  bool   `yaml:"shouldNotContainInterfaces"`
	ShouldNotContainTypes       bool   `yaml:"shouldNotContainTypes"`
	ShouldNotContainFunctions   bool   `yaml:"shouldNotContainFunctions"`
	ShouldNotContainMethods     bool   `yaml:"shouldNotContainMethods"`
}

type Config struct {
	DependenciesRules []DependenciesRule `yaml:"dependenciesRules"`
	ContentRules      []ContentsRule     `yaml:"contentsRules"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
