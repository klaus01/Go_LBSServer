package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/klaus01/Go_LBSServer/config"
	"github.com/klaus01/Go_LBSServer/controllers"
	"github.com/klaus01/Go_LBSServer/database"
)

func main() {
	environment := flag.String("e", "development", "")
	config.Init(*environment)
	database.Init()

	r := gin.Default()
	r.Static("/public", "./public")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	controllers.Init(r)
	r.Run()
}
