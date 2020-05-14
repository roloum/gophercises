package cmd

import (
	"fmt"
	"strings"

	"github.com/roloum/gophercises/task/internal/models"
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

		//Create task model
		//taskModel := models.NewDatastore(&db.Bolt{Name: "task.db"})
		ctx := cmd.Context()
		taskModel := ctx.Value(TaskModelKey).(*models.DataStore)

		task, err := taskModel.CreateTask(taskName)
		if err != nil {
			er(err)
		}
		fmt.Printf("Task created with ID: %d\n", task.ID)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	/*
		addCmd.Flags().StringVarP(&name, "name", "n", "", "task name")
		addCmd.MarkFlagRequired("name")
	*/
}
