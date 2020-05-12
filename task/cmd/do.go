package cmd

import (
	"fmt"
	"strconv"

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

		ids := []int{}
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error parsing: %s\n", arg)
			}
			ids = append(ids, id)
		}
		fmt.Println("Do tasks ", ids)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	/*
		doCmd.Flags().IntVarP(&taskID, "task-id", "t", 0, "Task ID")
		doCmd.MarkFlagRequired("task-id")
	*/
}
