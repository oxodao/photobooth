package cmd

import (
	"fmt"

	"github.com/oxodao/photobooth/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current photobooth version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--- Photobooth by Oxodao ---")
		fmt.Println("Current version: ", utils.CURRENT_VERSION)
		fmt.Println("Current commit: ", utils.CURRENT_COMMIT)
	},
}
