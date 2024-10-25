package main

import (
	"example.com/abobtech/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()
	route.POST("/ping", controller.CreateSchedule)
	route.GET("/getAllSchedules", controller.GetAllSchedules)
	// route.DELETE("/deleteSchedule/:priority", controller.DeleteSchedule)
	route.DELETE("/deleteSchedule", controller.DeleteSchedule)
	route.PATCH("/updateSchedule", controller.UpdateSchedule)

	route.Run("localhost:8080")
}
