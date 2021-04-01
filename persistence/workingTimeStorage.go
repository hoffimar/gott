package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hoffimar/gott/types"
)

func AddWorkingTime(writer io.Writer, inputInterval types.WorkingInterval) {

	var workingTimeList []types.WorkingInterval
	workingTimeList = append(workingTimeList, inputInterval)

	j, err := json.Marshal(workingTimeList)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintf(writer, string(j))
	if err != nil {
		fmt.Printf("Error printing to file: %s", err)
	}
}

func SaveWorkingTimes(writer io.Writer, workingTimeList []types.WorkingInterval) {
	j, err := json.Marshal(workingTimeList)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintf(writer, string(j))
	if err != nil {
		fmt.Printf("Error printing to file: %s", err)
	}
}

func GetWorkingTimes(filename string) (result []types.WorkingInterval, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}

	return result, nil
}
