package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use: "root",
	//Short: "Git is a distributed version control system.",
	Args: cobra.ExactArgs(0),
}
