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

	{
		v1 := route.Group("/api/v1/utility")
		v1.POST("/signup", controller.SignUp)
		v1.POST("/login", controller.LogIn)
		v1.GET("/getUsers", controller.GetUsers)
		v1.PATCH("/deposit", controller.Deposit)
		v1.PATCH("/debit", controller.Debit)
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	if err := route.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
