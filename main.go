package main

import (
	"rkenum/cmd"
)

var (
	BuildDate  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

func main() {
	cmd.SetVersionInfo(cmd.VersionInfo{
		Version:   GitTag,
		Commit:    CommitHash,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	})
	cmd.Execute()
}
