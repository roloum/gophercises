package cmd

import (
	"fmt"

	"github.com/roloum/gophercises/task/internal/models"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pending tasks",
	Run: func(cmd *cobra.Command, args []string) {

		/*
			//Create task model
			taskModel := models.NewDatastore(&db.Bolt{Name: "task.db"})
			fmt.Println("created model")
		*/
		ctx := cmd.Context()
		taskModel := ctx.Value(TaskModelKey).(*models.DataStore)

		taskList, err := taskModel.LoadTasks()
		if err != nil {
			er(err)
		}

		if len(taskList) == 0 {
			fmt.Println("Empty")
		}

		for _, task := range taskList {
			fmt.Printf("%d: %s\n", task.ID, task.Name)
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
