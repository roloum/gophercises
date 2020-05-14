package db

import (
	"strconv"

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
func (b *Bolt) bucket() error {

	return b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(taskBucket)); err != nil {
			return err
		}
		return nil
	})
}

//CreateTask creates a task in the bucket
func (b *Bolt) CreateTask(t *models.Task) error {

	if err := b.bucket(); err != nil {
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

	if err := b.bucket(); err != nil {
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
