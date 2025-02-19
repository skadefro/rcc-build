package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var (
	configFile   string
	profileName  string
	clearProfile bool
)

var configureExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a configuration profile for Robocorp tooling.",
	Long:  "Export a configuration profile for Robocorp tooling.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Configuration export lasted").Report()
		}
		profile := loadNamedProfile(profileName)
		err := profile.SaveAs(configFile)
		pretty.Guard(err == nil, 1, "Error while exporting profile, reason: %v", err)
		pretty.Ok()
	},
}

func init() {
	configureCmd.AddCommand(configureExportCmd)
	configureExportCmd.Flags().StringVarP(&configFile, "filename", "f", "exported_profile.yaml", "The filename where configuration profile is exported.")
	configureExportCmd.Flags().StringVarP(&profileName, "profile", "p", "", "The name of configuration profile to export.")
	configureExportCmd.MarkFlagRequired("profile")
}
