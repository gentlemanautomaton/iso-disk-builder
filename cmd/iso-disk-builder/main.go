package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var cli struct {
		Build BuildCmd `kong:"cmd,help='Builds an ISO 9660 disk image from a source directory.'"`
	}

	parser := kong.Must(&cli,
		kong.Description("Builds ISO 9660 disk images."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	app, parseErr := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(parseErr)

	appErr := app.Run()
	app.FatalIfErrorf(appErr)
}
