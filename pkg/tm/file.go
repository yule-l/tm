package tm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const resultFilePattern = "%s.done"
const defaultPerm = 0755

type tasksFileManager struct {
	workDir string
}

func newTasksFileManager() (*tasksFileManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &tasksFileManager{
		workDir: filepath.Join(homeDir, ".yule-l", "tm"),
	}, nil
}

func (t tasksFileManager) createWorkDirIfNoExists() {
	_ = os.MkdirAll(t.workDir, defaultPerm)
}

func (t tasksFileManager) setup(cfg Config) (done io.ReadWriter, tasks *tasks, queue *bufio.Scanner, err error) {
	t.createWorkDirIfNoExists()

	done, err = t.openDoneFile(filepath.Join(t.workDir, t.getDoneFilePath(cfg.FilePath)), cfg.Force)
	if err != nil {
		return
	}
	tasks = t.loadTasksInfo(done)
	queue, err = t.openQueueFile(cfg.FilePath)
	return
}

func (t tasksFileManager) getDoneFilePath(filePath string) string {
	return fmt.Sprintf(resultFilePattern, filepath.Base(filePath))
}

func (t tasksFileManager) openDoneFile(filePath string, remove bool) (io.ReadWriter, error) {
	if remove {
		_ = os.Truncate(filePath, 0)
	}
	return os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, defaultPerm)
}

func (tasksFileManager) loadTasksInfo(file io.Reader) *tasks {
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

func (tasksFileManager) openQueueFile(filePath string) (*bufio.Scanner, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, defaultPerm)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return scanner, nil
}
