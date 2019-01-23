package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time to the version of azssh being built
// Set it via `go build -ldflags "-X github.com/noelbundick/azssh/cmd.Version=$VERSION"`
var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of azssh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
