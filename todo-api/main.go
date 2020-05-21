package main

import (
	"github.com/layemut/todo-application/todo-api/app"
	"github.com/layemut/todo-application/todo-api/config"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.Initialize(config)
	app.Run(":8060")
}
