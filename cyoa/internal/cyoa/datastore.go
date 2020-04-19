package cyoa

//StoryAccessObject ...
type StoryAccessObject interface {
	LoadStory(storyID string, story *Story) error
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
func (d *DataStore) LoadStory(storyID string) (Story, error) {
	story := Story{}
	err := d.dao.LoadStory(storyID, &story)
	return story, err
}
