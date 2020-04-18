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

//LoadStory ...
func (j *JSON) LoadStory(fileName string) (cyoa.Story, error) {

	if fileName == "" {
		return nil, errors.New("Story is empty")
	}

	j.Log.Printf("Loading story: %v\n", fileName)

	return parseJSONFile(fileName, j.Log)
}

func parseJSONFile(fileName string, log *log.Logger) (cyoa.Story, error) {

	log.Printf("Opening file: %v\n", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	var story cyoa.Story

	log.Printf("Decoding file: %v\n", fileName)

	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
