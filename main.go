/*
Copyright © 2022 Lammaskoira authors

*/
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/lammaskoira/bark/cmd"
)

func runCommand() int {
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := cmd.Execute(ctx); err != nil {
		return 1
	}

	return 0
}

func main() {
	os.Exit(runCommand())
}
