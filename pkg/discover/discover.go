package discover

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/distributed-technologies/flux-graph/pkg/helmRelease"
	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
)

// Gets a list of *.yaml files that contains `apiVersion: argocd-discover/v1alpha1` string and generates an ArgoCD application resource that is written to stdout
func Discover(root string, folder string, showHelmReleases bool) error {
	yamlFiles, err := GetFilesThatContains(filepath.Join(root, folder), "apiVersion: kustomize.toolkit.fluxcd.io")
	if err != nil {
		return err
	}

	for _, path := range yamlFiles {

		var ks kustomization.Kustomization

		ks.GetValuesFromYamlFile(path)

		if ks.HasDependsOn() {
			helmFiles, err := GetFilesThatContains(filepath.Join(root, ks.Spec.Path), "helm.toolkit.fluxcd.io")
			if err != nil {
				return err
			}

			for _, hrPath := range helmFiles {

				var hr helmRelease.HelmRelease
				hr.GetValuesFromYamlFile(hrPath)

				if hr.HasDependsOn() && showHelmReleases {
					hr.Parent = ks.Metadata.Name
					ks.HRSlice = append(ks.HRSlice, hr.Metadata.Name)
					helmRelease.HelmReleases = append(helmRelease.HelmReleases, hr)
				}

			}

			kustomization.Kustomizations = append(kustomization.Kustomizations, ks)

		}
	}

	return nil
}

// Looks up all files in the root, and checks if it contains the 'argocd-discover' apiVersion
func GetFilesThatContains(folder string, contains string) ([]string, error) {
	yamlFiles := []string{}
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".yaml") {

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			stringFile := string(file)

			splitFiles := strings.Split(stringFile, "---")

			for _, content := range splitFiles {
				if strings.Contains(content, contains) {
					tmpFile, err := os.CreateTemp(os.TempDir(), "*.yaml")
					if err != nil {
						return err
					}
					tmpFile.WriteString(content)
					yamlFiles = append(yamlFiles, tmpFile.Name())
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return yamlFiles, nil
}
