package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/roloum/gophercises/task/internal/models"
	"github.com/spf13/cobra"
)

//var taskID int

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Runs a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("do called with taskID: %d\n", taskID)

		ids := getIDs(args)
		if len(ids) > 0 {
			ctx := cmd.Context()
			taskModel := ctx.Value(TaskModelKey).(*models.DataStore)

			performed, err := taskModel.DoTasks(ids)
			if len(performed) > 0 {
				fmt.Printf("Performed:\n%s\n", strings.Join(performed, "\n"))
			}
			if err != nil {
				fmt.Println("Errors:")
				er(err)
			}
		}
	},
}

func getIDs(args []string) []int {
	ids := []int{}
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Error parsing: %s\n", arg)
		} else {
			ids = append(ids, id)
		}
	}
	return ids
}

func init() {
	RootCmd.AddCommand(doCmd)

	/*
		doCmd.Flags().IntVarP(&taskID, "task-id", "t", 0, "Task ID")
		doCmd.MarkFlagRequired("task-id")
	*/
}
