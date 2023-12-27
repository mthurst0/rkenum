package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type VersionInfo struct {
	Version   string
	Commit    string
	BuildDate string
	GoVersion string
}

var Version VersionInfo

func SetVersionInfo(v VersionInfo) {
	Version = v
}

func createVersionCmd() *cobra.Command {
	var enumCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version    : %s\n", Version.Version)
			fmt.Printf("Commit     : %s\n", Version.Commit)
			fmt.Printf("Build date : %s\n", Version.BuildDate)
			fmt.Printf("Go version : %s\n", Version.GoVersion)
			return nil
		},
	}
	return enumCmd
}

func init() {
	rootCmd.AddCommand(createVersionCmd())
}
