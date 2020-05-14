package models

//TaskAccessObject interface defines all methods need to be implemented
//by the database connector
type TaskAccessObject interface {
	LoadTasks() ([]Task, error)
	CreateTask(t *Task) error
	DoTasks(taskIDs []int) ([]string, error)
}

//DataStore contains a data access object that implements the
//TaskAccessObject interface
type DataStore struct {
	dao TaskAccessObject
}

//NewDatastore returns a new datastore object
func NewDatastore(dao TaskAccessObject) *DataStore {
	d := DataStore{dao}
	return &d
}

//LoadTasks returns a list of Tasks
func (d *DataStore) LoadTasks() ([]Task, error) {
	return d.dao.LoadTasks()
}

//CreateTask creates a task
func (d *DataStore) CreateTask(name string) (*Task, error) {
	t := &Task{Name: name}

	err := d.dao.CreateTask(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

//DoTasks performs all the tasks in the taskIDs array
func (d *DataStore) DoTasks(taskIDs []int) ([]string, error) {
	return d.dao.DoTasks(taskIDs)
}
