package parkings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type JsonStateRepository struct {
	directory string
}

func ProvideJsonStateRepository(directory string) JsonStateRepository {
	if directory == "" {
		directory = "./data"
	}

	os.MkdirAll(directory, os.ModePerm)

	return JsonStateRepository{directory}
}

func (sr *JsonStateRepository) SaveAll(states []State, currentTime time.Time) error {
	data, _ := json.MarshalIndent(states, "", "    ")

	return ioutil.WriteFile(sr.getFileName(currentTime), data, 0644)
}

func (sr *JsonStateRepository) getFileName(currentTime time.Time) string {
	return sr.directory + "/" + currentTime.Format("2006-01-02_15-04") + ".json"
}
