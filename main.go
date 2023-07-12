package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mrnbort/webhook.git/api"
	"github.com/mrnbort/webhook.git/executor"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type options struct {
	Port     string `short:"p" long:"port" description:"port" default:":8080"`
	FileName string `short:"f" long:"file" description:"file to read task commands from" default:"config.yaml"`
}

func main() {
	var opts options
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v", err)
		}
		os.Exit(1)
	}

	if err := run(opts); err != nil {
		log.Panicf("[ERROR] %v", err)
	}
}

func run(opts options) error {
	data, err := executor.NewTaskFinder(opts.FileName)
	if err != nil {
		log.Panicf("[ERROR] failed to read configuration file: %v", err)
	}

	apiService := api.Service{
		Executor: data,
		Port:     opts.Port,
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // cancel on SIGINT or SIGTERM
	go func() {
		sig := <-sigs
		log.Printf("received signal: %v", sig)
		cancel()
	}()

	if err := apiService.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("webhook service canceled")
			return nil
		}
		return fmt.Errorf("webhook service failed: %v", err)
	}
	return nil
}
