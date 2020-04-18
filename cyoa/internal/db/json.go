package db

import (
	"log"

	"github.com/roloum/gophercises/cyoa/internal/cyoa"
)

//JSON ...
type JSON struct {
	Dir  string
	File string
	Log  *log.Logger
}

//LoadStory ...
func (j *JSON) LoadStory(story interface{}) (cyoa.Story, error) {

	log.Printf("Loading story: %v\n", story.(string))
	return cyoa.Story{}, nil
}
