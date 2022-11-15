/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/distributed-technologies/flux-graph/pkg/discover"
	"github.com/distributed-technologies/flux-graph/pkg/graph"
	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
	"github.com/distributed-technologies/flux-graph/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const discoverDesc = `
Builds a DOT graph based in files that has the 'apiVersion: kustomize.toolkit.fluxcd.io' and has a 'dependsOn field'.
`

// buildCmd represents the discover command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Discovers any files, in the current folder or subfolder, that contains 'apiVersion: kustomize.toolkit.fluxcd.io' and generates a dependency graph",
	Long:  discoverDesc,
	Run: func(cmd *cobra.Command, args []string) {
		folder := viper.GetString("folder")
		root := viper.GetString("root-folder")
		helmReleases := viper.GetBool("show-helm-releases")

		logging.Debug("folder: %v\n", folder)
		logging.Debug("root: %v\n", root)
		logging.Debug("helmReleases: %v\n", helmReleases)

		err := discover.Discover(root, folder, helmReleases)
		if err != nil {
			panic(err)
		}

		graphString, err := graph.New("main").Generate(kustomization.Kustomizations)
		if err != nil {
			panic(err)
		}

		fmt.Printf(graphString)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().String("folder", "./", "Folder (relative to the root-folder) to find apps in (recursive)")
	buildCmd.Flags().String("root-folder", "./", "The root of the flux folder")
	buildCmd.Flags().Bool("show-helm-releases", false, "Shows in line helm release dependencies")

	viper.BindPFlags(buildCmd.Flags())
}
