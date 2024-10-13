package main

import (
	"example.com/abobtech/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()
	route.GET("/ping", controller.CreateTask)

	route.Run("localhost:8080")
}
