/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/distributed-technologies/flux-graph/pkg/discover"
	graphvizWrapper "github.com/distributed-technologies/flux-graph/pkg/graphviz-warpper"
	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
	"github.com/distributed-technologies/flux-graph/pkg/logging"
	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const discoverDesc = `
The discover cmd builds ArgoCD application resources based on any file, it finds in the current folder and subfolders,
that contains 'apiVersion: argocd-discover/v1alpha1' as the first line of the file.

This is done by walking through the folder structure reading the first line of any '*.yaml' file,
and checking if it matches 'apiVersion: argocd-discover/v1alpha1' if it does, it reads the rest of the content,
which describes the chart that should be used and where it should be deployed.
`

// buildCmd represents the discover command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Discovers any files, in the current folder or subfolder, that contains 'apiVersion: argocd-discover/v1alpha1'.",
	Long:  discoverDesc,
	Run: func(cmd *cobra.Command, args []string) {
		folder := viper.GetString("folder")
		outDir := viper.GetString("outDir")

		err := discover.Discover(folder)
		if err != nil {
			panic(err)
		}

		g := graphviz.New()
		graph, err := g.Graph()
		if err != nil {
			log.Fatal(err)
		}
		graph.SetDPI(256)

		defer func() {
			if err := graph.Close(); err != nil {
				log.Fatal(err)
			}
			g.Close()
		}()

		gw := graphvizWrapper.GraphWrap{
			G:  graph,
			Gv: g,
		}

		for _, v := range kustomization.Kustomizations {
			logging.Debug("ks: %v", v.Name())
			logging.Debug("Dependencies: %v\n", v.GetDepndencies())

			if n, _ := gw.G.Node(v.Name()); n == nil {
				gw.MakeNode(v.Name())
			}

		}

		for _, ks := range kustomization.Kustomizations {

			startN, err := gw.G.Node(ks.Name())
			if err != nil {
				logging.WrapError("err: %e", err)
			}

			for _, dep := range ks.GetDepndencies() {
				logging.Debug("dep: %s", dep)

				depN, _ := gw.G.Node(dep)
				if depN == nil {
					depN = gw.MakeNode(dep)
				}

				if err != nil {
					logging.WrapError("err: %e", err)
				}
				gw.MakeEdge(depN, startN)
			}
		}
		gw.Render(outDir)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().String("folder", "./", "Folder to find apps in (recursive)")
	buildCmd.Flags().String("outDir", "./", "Folder to output the graph")

	viper.BindPFlags(buildCmd.Flags())
}
