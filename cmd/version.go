package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Godigger",
  Long:  `All software has versions. This is Godigger's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Godigger Microservice Generator v0.1 -- HEAD")
  },
}