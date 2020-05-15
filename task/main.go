/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/roloum/gophercises/task/cmd"
	"github.com/roloum/gophercises/task/internal/db"
	"github.com/roloum/gophercises/task/internal/models"
)

func main() {

	/*
		Testing the database timeout
		go func() {
			crap := db.Bolt{Name: "task.db"}
			if err := crap.Connect(); err != nil {
				er(err)
			}
		}()
	*/

	time.Sleep(1 * time.Second)

	//Setup database connection
	dao := db.Bolt{Name: "task.db"}
	if err := dao.Connect(db.WithTimeout(1 * time.Second)); err != nil {
		er(err)
	}
	defer dao.Close()

	//Create task model
	taskModel := models.NewDatastore(&dao)
	//Store model in context to pass to command
	ctx := context.WithValue(nil, cmd.TaskModelKey, taskModel)

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		er(err)
	}

}

func er(err error) {
	fmt.Println(err)
	os.Exit(1)
}
