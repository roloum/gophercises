package cyoa

//StoryAccessObject ...
type StoryAccessObject interface {
	LoadStory(story string) (Story, error)
}

//DataStore ...
type DataStore struct {
	dao StoryAccessObject
}

//NewDataStore ...
func NewDataStore(dao StoryAccessObject) *DataStore {
	d := DataStore{dao}
	return &d
}

//LoadStory ...
func (d *DataStore) LoadStory(story string) (Story, error) {
	return d.dao.LoadStory(story)
}
