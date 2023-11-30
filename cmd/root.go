package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isDebug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flux-graph",
	Short: "Builds a dependency graph from a flux repo",
	Long: `Flux-graph is a tool that looks at kustomization files and builds a dependency graph is any exists.`,

	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent Flags will be available to this command and all subcommands to this
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "enable debug logs")

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("FG") // Standing for 'flux-graph'
	viper.AutomaticEnv()     // read in environment variables that match

}
