package main

import (
	"log"
	"os"

	"example.com/abobtech/controller"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	route := gin.Default()
	// rateLimit := controller.FixedWindowRateLimiter(5, time.Minute)
	route.POST("/signup", controller.SignUp)
	route.POST("/login", controller.LogIn)
	route.GET("/getUsers", controller.GetUsers)
	route.PATCH("/deposit", controller.Deposit)
	route.PATCH("/debit", controller.Debit)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := route.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
