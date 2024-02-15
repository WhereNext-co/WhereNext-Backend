package main

import (
	"github.com/WhereNext-co/WhereNext-Backend.git/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
 initializers.LoadEnvVariables()
	}


func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}