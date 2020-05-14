package db

import (
	"encoding/binary"
	"fmt"
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
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		bu := tx.Bucket([]byte(taskBucket))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := bu.NextSequence()
		t.ID = int(id)

		// Persist bytes
		return bu.Put(itob(t.ID), []byte(t.Name))
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
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
			fmt.Printf("k: %s, v: %s\n", k, v)
			id, _ := strconv.Atoi(string(k))

			list = append(list, models.Task{ID: id, Name: string(v)})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return list, nil
}
