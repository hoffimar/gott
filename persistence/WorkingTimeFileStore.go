package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/hoffimar/gott/types"
)

type WorkingTimeFileStore struct {
	filePath string
}

func NewWorkingTimeFileStore(dirName string, fileName string) (store *WorkingTimeFileStore, err error) {
	os.MkdirAll(dirName, 0700)

	filePath := path.Join(dirName, fileName)

	// Check for file existence
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Creating file %s", filePath)

		// fill file with initial (empty) array
		content := []byte("[]")
		ioutil.WriteFile(filePath, content, 0600)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	return &WorkingTimeFileStore{filePath: filePath}, nil
}

func (fileStore WorkingTimeFileStore) AddWorkingTime(inputInterval types.WorkingInterval) (err error) {

	workingTimes, _ := fileStore.GetWorkingTimes()
	workingTimes = append(workingTimes, inputInterval)

	bytes, err := json.Marshal(workingTimes)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileStore.filePath, bytes, 0600)
}

func (fileStore WorkingTimeFileStore) GetWorkingTimes() (result []types.WorkingInterval, err error) {
	data, err := ioutil.ReadFile(fileStore.filePath)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}
