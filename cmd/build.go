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

		logging.Debug("folder: %s\n", folder)

		err := discover.Discover(folder)
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

	buildCmd.Flags().String("folder", "./", "Folder to find apps in (recursive)")

	viper.BindPFlags(buildCmd.Flags())
}
