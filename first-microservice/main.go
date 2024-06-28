package main

import (
	"context"
	"fmt"

	"github.com/F-Madruga/go-basics/application"
)

func main() {
	app := application.NewApp()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app", err)
	}
}
