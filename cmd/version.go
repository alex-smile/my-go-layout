package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"mygo/template/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		info := []string{
			"Version: " + version.Version,
			"Commit: " + version.Commit,
			"Build Time: " + version.BuildTime,
			"Go Version: " + version.GoVersion,
		}
		fmt.Println(strings.Join(info, "\n"))
	},
}
