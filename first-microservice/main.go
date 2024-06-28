package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/F-Madruga/go-basics/application"
)

func main() {
	app := application.NewApp()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
