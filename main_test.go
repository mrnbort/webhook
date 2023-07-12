package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	opts := options{
		Port:     ":8081",
		FileName: "config.yaml",
	}

	// Run the run function with the test options in a separate goroutine
	go func() {
		err := run(opts)
		assert.NoError(t, err)
	}()

	// Wait for 5 seconds to let the run function run and then cancel it
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	signal.NotifyContext(ctx, os.Interrupt)
}
