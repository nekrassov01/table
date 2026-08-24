package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nekrassov01/table/internal/skills"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := skills.RunGoldens(ctx, os.Args[1:], os.Stdout, os.Stderr)
	cancel()
	os.Exit(code)
}
