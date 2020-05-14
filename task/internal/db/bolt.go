package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/roloum/gophercises/task/internal/models"
)

/*
db, err := bolt.Open("task.db", 0600, nil)
if err != nil {
  er(err)
}
*/

const fileMode = 0600
const taskBucket = "task"
const taskDoneBucket = "done"

//Bolt database adapter implements TaskAccessObject interface
type Bolt struct {
	Name    string
	Options *bolt.Options
	db      *bolt.DB
}

//Connect establishes a connection to the Bolt database
func (b *Bolt) Connect() (err error) {

	b.db, err = bolt.Open(b.Name, fileMode, b.Options)
	if err != nil {
		return err
	}
	return
}

//Close database connection
func (b *Bolt) Close() {
	b.db.Close()
}

//bucket Creates the bucket if it does not exist
func (b *Bolt) bucket(bucketName string) error {

	return b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return err
		}
		return nil
	})
}

//DoTasks runs the list of tasks provided on taskIDs using the following steps
//for each of the Tasks
// - begin transaction
// - load task
// - insert task into done bucket
// - delete task from task bucket
// - commit transaction
func (b *Bolt) DoTasks(taskIDs []int) ([]string, error) {
	errorDescription, performed := []string{}, []string{}

	//Make sure buckets exist
	if err := b.bucket(taskBucket); err != nil {
		return nil, err
	}
	if err := b.bucket(taskDoneBucket); err != nil {
		return nil, err
	}

	for _, taskID := range taskIDs {

		err := func(taskID int) error {

			//begin
			tx, err := b.db.Begin(true)
			if err != nil {
				return err
			}
			defer tx.Rollback()

			taskIDb := []byte(strconv.Itoa(taskID))
			//load task
			tb := tx.Bucket([]byte(taskBucket))
			taskName := tb.Get(taskIDb)
			if taskName == nil {
				return fmt.Errorf("Task ID %d does not exist", taskID)
			}

			//create done record
			db := tx.Bucket([]byte(taskDoneBucket))
			if err = db.Put(taskIDb, taskName); err != nil {
				return err
			}

			//delete task record
			if err = tb.Delete(taskIDb); err != nil {
				return err
			}

			//commit
			if err := tx.Commit(); err != nil {
				return err
			}

			performed = append(performed, string(taskName))

			return nil
		}(taskID)
		if err != nil {
			errorDescription = append(errorDescription, err.Error())
		}
	}

	var err error
	if len(errorDescription) > 0 {
		err = errors.New(strings.Join(errorDescription, "\n"))
	}
	return performed, err
}

//CreateTask creates a task in the bucket
func (b *Bolt) CreateTask(t *models.Task) error {

	if err := b.bucket(taskBucket); err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte(taskBucket))

		id, _ := bu.NextSequence()
		t.ID = int(id)

		return bu.Put([]byte(strconv.Itoa(t.ID)), []byte(t.Name))
	})
}

//LoadTasks loads all the tasks into a TaskList struct
func (b *Bolt) LoadTasks() ([]models.Task, error) {
	list := []models.Task{}

	if err := b.bucket(taskBucket); err != nil {
		return nil, err
	}

	if err := b.db.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte(taskBucket)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			id, _ := strconv.Atoi(string(k))

			list = append(list, models.Task{ID: id, Name: string(v)})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return list, nil
}
