package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//TaskModelKey ...
const TaskModelKey = "taskModel"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "task",
	Short: "A CLI Task Manager",
}

func er(err error) {
	fmt.Println(err)
	os.Exit(1)
}
