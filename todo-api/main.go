package main

import (
	"fmt"

	"github.com/layemut/todo-api/app"
	"github.com/layemut/todo-api/config"
)

func main() {
	config := config.GetConfig()
	fmt.Println(config.DB.Host)
	app := &app.App{}
	app.Initialize(config)
	app.Run(":8080")
}
