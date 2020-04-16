package parkings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type JsonStateRepository struct {
	directory string
}

func ProvideJsonStateRepository(directory string) JsonStateRepository {
	if directory == "" {
		directory = "./data"
	}

	return JsonStateRepository{directory}
}

func (sr *JsonStateRepository) SaveAll(states []State, currentTime time.Time) error {
	data, _ := json.MarshalIndent(states, "", "    ")
	filename := sr.getFileName(currentTime)

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (sr *JsonStateRepository) getFileName(currentTime time.Time) string {
	return sr.directory + "/" + currentTime.Format("2006-01-02/15-04") + ".json"
}
