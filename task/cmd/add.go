package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

//var name string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := strings.Join(args, " ")

		fmt.Println("Add task: ", taskName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	/*
		addCmd.Flags().StringVarP(&name, "name", "n", "", "task name")
		addCmd.MarkFlagRequired("name")
	*/
}
