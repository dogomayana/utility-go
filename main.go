package main

import (
	"log"
	"os"

	"example.com/abobtech/controller"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

func main() {

	route := gin.Default()
	route.POST("/createSchedule", controller.CreateSchedule)
	route.GET("/getAllSchedules", controller.GetAllSchedules)
	route.DELETE("/deleteSchedule", controller.DeleteSchedule)
	route.PATCH("/updateSchedule", controller.UpdateSchedule)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := route.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

	// route.Run("localhost:8080")
}
