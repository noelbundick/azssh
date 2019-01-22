package cmd

import (
	"fmt"
	"os"

	"github.com/noelbundick/azssh/pkg/azssh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azssh",
	Short: "Launch Azure Cloud Shell from your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		token := azssh.GetToken()
		url := azssh.ProvisionCloudShell(token, shellType)
		azssh.ConnectToWebsocket(url)
	},
}
var shellType string

func init() {
	rootCmd.Flags().StringVarP(&shellType, "shell", "s", "bash", "shell to launch (bash / pwsh)")
}

// Execute launches a Cloud Shell and connects it to STDIN/STDOUT
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
