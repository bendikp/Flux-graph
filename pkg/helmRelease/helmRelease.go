package helmRelease

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Structs for handling the information that comes from the helmRelease api's
type HelmRelease struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec"`
	Parent   string
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Spec struct {
	DependsOn []DependsOn `yaml:"dependsOn"`
}

type DependsOn struct {
	Name string `yaml:"name"`
}

var HelmReleases []HelmRelease

// Unmarshals a .yaml file into the app struct
func (hr *HelmRelease) GetValuesFromYamlFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, &hr)
}

func (hr *HelmRelease) HasDependsOn() bool {
	return hr.Spec.DependsOn != nil
}

// Enclosing name in double quotes since the graph language wants them to be in quotes
func (hr *HelmRelease) Name() string {
	return "\"" + hr.Metadata.Name + "\""
}

func (hr *HelmRelease) GetDependencies() []string {
	dependencies := []string{}

	for _, v := range hr.Spec.DependsOn {

		// Enclosing name in double quotes since the graph language wants them to be in quotes
		dependencies = append(dependencies, "\""+v.Name+"\"")
	}
	return dependencies
}
