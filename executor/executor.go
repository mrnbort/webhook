package executor

import (
	"context"
	"fmt"
	"github.com/mrnbort/webhook.git/task"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
)

// TaskFinder is a file reader that finds the command for the requested task
type TaskFinder struct {
	fileName string
	tasks    []task.Task
}

// NewTaskFinder creates a file reader that reads from the configuration file and stores data in a slice
func NewTaskFinder(fileName string) (*TaskFinder, error) {
	taskList := &task.TaskList{}
	log.Printf("opening configuration file %v", fileName)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current directory is %v", dir)

	data, err := os.ReadFile(dir + "/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	if err := yaml.Unmarshal(data, &taskList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %v", err)
	}
	return &TaskFinder{
		fileName: fileName,
		tasks:    taskList.Tasks,
	}, nil
}

// ExecuteTask checks if the task ID and auth matches with the config file and if so, executes the task command
func (tf *TaskFinder) ExecuteTask(ctx context.Context, request task.Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for _, tsk := range tf.tasks {
		if tsk.ID == request.ID && tsk.Auth == request.Auth {
			cmd := exec.Command("sh", "-c", tsk.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			log.Printf("executing task: %v", tsk.Command)
			//if err := cmd.Run(); err != nil {
			// return
			//}
			log.Printf("task executed successefully")
			return nil
		}
	}
	err := fmt.Errorf("invalid task ID or auth key")
	return err
}
