package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/mrnbort/webhook.git/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"
)

func TestService_Run(t *testing.T) {
	done := make(chan struct{})
	go func() {
		<-done
		e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.NoError(t, e)
	}()
}

func TestService_postExecuteTask(t *testing.T) {
	exectr := &ExecutorMock{
		ExecuteTaskFunc: func(ctx context.Context, request task.Task) error {
			return nil
		},
	}
	svc := &Service{
		Executor: exectr,
	}

	ts := httptest.NewServer(svc.routes())
	defer ts.Close()

	client := http.Client{Timeout: time.Second}

	{ // successful attempt
		url := fmt.Sprintf("%s/execute?id=test&key=123", ts.URL)
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"status":"ok"}`+"\n", string(data))
		require.Equal(t, 1, len(exectr.ExecuteTaskCalls()))
		assert.Equal(t, "test", exectr.ExecuteTaskCalls()[0].Request.ID)
		assert.Equal(t, "123", exectr.ExecuteTaskCalls()[0].Request.Auth)
	}

	{ // failed post
		exectr.ExecuteTaskFunc = func(ctx context.Context, request task.Task) error {
			return errors.New("oh oh")
		}
		req, err := http.NewRequest("POST", ts.URL+"/execute?id=test&key=123", nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"error":"oh oh"}`+"\n", string(data))
		require.Equal(t, 2, len(exectr.ExecuteTaskCalls()))
	}
}
