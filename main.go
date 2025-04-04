package main

import (
	"log"
	"os"
	"time"

	"example.com/abobtech/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	rateLimit := controller.FixedWindowRateLimiter(5, time.Minute)
	// route.POST("/createSchedule", controller.CreateSchedule)
	route.POST("/create", rateLimit, controller.CreateItem)

	route.GET("/getAllSchedules", controller.GetAllSchedules)
	route.GET("/getSchedule", controller.GetSchedule)
	// route.DELETE("/deleteSchedule", controller.DeleteSchedule)
	// route.PATCH("/updateSchedule", controller.UpdateSchedule)
	port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }
	if err := route.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
