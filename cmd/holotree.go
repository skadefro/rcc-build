package cmd

import (
	"github.com/robocorp/rcc/settings"
	"github.com/spf13/cobra"
)

var holotreeCmd = &cobra.Command{
	Use:     "holotree",
	Aliases: []string{"ht"},
	Short:   "Group of holotree commands.",
	Long:    "Group of holotree commands.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		settings.CriticalEnvironmentSettingsCheck()
	},
}

func init() {
	rootCmd.AddCommand(holotreeCmd)
}
