package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klaus01/Go_LBSServer/config"
	"github.com/klaus01/Go_LBSServer/database"
)

func main() {
	environment := flag.String("e", "development", "")
	config.Init(*environment)
	database.Init()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "朝秦暮楚",
		})
	})
	r.Run()
}
