package adventure

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func Decode(fileName string) (*Story, error) {
	var path string
	var err error
	path, err = filepath.Abs(fileName)
	storyData, err := ioutil.ReadFile(path)
	var story *Story
	err = json.Unmarshal(storyData, &story)
	return story, err
}
