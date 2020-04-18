package cyoa

//StoryAccessObject ...
type StoryAccessObject interface {
	LoadStory(story interface{}) (Story, error)
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
func (d *DataStore) LoadStory(story interface{}) (Story, error) {
	return d.dao.LoadStory(story)
}

//Memeo ...
func Memeo() {

}

//FFFFFFFFF ...
type FFFFFFFFF struct{}
