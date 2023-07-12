package executor

import (
	"context"
	"fmt"
	"github.com/mrnbort/webhook.git/task"
	"testing"
)

func TestTaskFinder_ExecuteTask(t *testing.T) {
	tasks := []task.Task{{ID: "test1", Auth: "1234", Command: "echo hello"}}
	tf := &TaskFinder{tasks: tasks}

	// Test that a valid request executes without error
	req := task.Task{ID: "test1", Auth: "1234"}
	if err := tf.ExecuteTask(context.Background(), req); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test that a request with an invalid task ID returns an error
	req = task.Task{ID: "invalid", Auth: "1234"}
	wantErr := fmt.Errorf("invalid task ID or auth key")
	if err := tf.ExecuteTask(context.Background(), req); err == nil || err.Error() != wantErr.Error() {
		t.Errorf("expected error %v; got %v", wantErr, err)
	}

	// Test that a request with an invalid task auth returns an error
	req = task.Task{ID: "test1", Auth: "kkk"}
	wantErr = fmt.Errorf("invalid task ID or auth key")
	if err := tf.ExecuteTask(context.Background(), req); err == nil || err.Error() != wantErr.Error() {
		t.Errorf("expected error %v; got %v", wantErr, err)
	}
}
