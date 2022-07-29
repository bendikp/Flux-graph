package discover

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
	"github.com/distributed-technologies/flux-graph/pkg/logging"
)

// Gets a list of *.yaml files that contains `apiVersion: argocd-discover/v1alpha1` string and generates an ArgoCD application resource that is written to stdout
func Discover(folder string) error {
	logging.Debug("folder: %v\n", folder)

	yamlFiles, err := GetFiles(folder)
	if err != nil {
		return err
	}

	for _, path := range yamlFiles {
		var ks kustomization.Kustomization

		ks.GetValuesFromYamlFile(path)

		if ks.HasDependsOn() {
			kustomization.Kustomizations = append(kustomization.Kustomizations, ks)
		}
	}

	return nil
}

// Looks up all files in the root, and checks if it contains the 'argocd-discover' apiVersion
func GetFiles(folder string) ([]string, error) {
	yamlFiles := []string{}
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".yaml") {

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "apiVersion: kustomize.toolkit.fluxcd.io") {
					yamlFiles = append(yamlFiles, path)
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	logging.Debug("yamlFiles: %s", yamlFiles)
	return yamlFiles, nil
}
