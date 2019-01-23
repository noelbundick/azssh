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

		resize := make(chan azssh.TerminalSize)
		initialSize := azssh.GetTerminalSize()

		url := azssh.ProvisionCloudShell(token, shellType, initialSize, resize)
		azssh.ConnectToWebsocket(url, resize)
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
