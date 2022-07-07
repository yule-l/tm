package tm

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

const resultFileSuffix = ".done"
const defaultPerm = 0755

func openDoneFile(filePath string, remove bool) (io.ReadWriter, error) {
	doneFilePathBuilder := strings.Builder{}
	doneFilePathBuilder.WriteString(filePath)
	doneFilePathBuilder.WriteString(resultFileSuffix)
	filePath = doneFilePathBuilder.String()
	if remove {
		_ = os.Truncate(filePath, 0)
	}
	return os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, defaultPerm)
}

func loadTasksInfo(file io.Reader) *tasks {
	t := newTasks()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		r := result{}
		err := json.Unmarshal(scanner.Bytes(), &r)
		if err != nil {
			continue
		}
		t.store(r.Payload, r)
	}
	return t
}

func openQueueFile(filePath string) (*bufio.Scanner, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, defaultPerm)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return scanner, nil
}
