package db

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/roloum/gophercises/cyoa/internal/cyoa"
)

//JSON ...
type JSON struct {
	Dir  string
	File string
	Log  *log.Logger
}

//LoadStories ...
func (j *JSON) LoadStories(dir string, stories *cyoa.Stories) {

}

//LoadStory ...
func (j *JSON) LoadStory(fileName string, story *cyoa.Story) error {

	if fileName == "" {
		return errors.New("Story is empty")
	}

	j.Log.Printf("Loading story: %v\n", fileName)

	return parseJSONFile(fileName, story, j.Log)
}

func parseJSONFile(fileName string, story *cyoa.Story, log *log.Logger) error {

	log.Printf("Opening file: %v\n", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	log.Printf("Decoding file: %v\n", fileName)

	if err := decoder.Decode(&story); err != nil {
		return err
	}

	return nil
}
