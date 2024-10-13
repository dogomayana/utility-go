package controller

import (
	"fmt"
	"net/http"

	"os"

	"example.com/abobtech/models"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

	"github.com/gin-gonic/gin"
)

func ConnectDB() (client *supabase.Client) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	client, err = supabase.NewClient(supabaseURL, supabaseAnonKey, nil)

	if err != nil {
		fmt.Println("cannot initalize client", err)
		return
	}
	return client
}

func CreateTask(c *gin.Context) {

	supaClient := ConnectDB()

	var json models.CreateSchedule
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type error"})
		return
	}
	row := models.CreateSchedule{
		Description: json.Description,
		Daymonth:    json.Daymonth,
		Priority:    json.Priority,
	}
	_, _, err := supaClient.From("scheduler").Insert(row, false, "", "", "").Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server Error",
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"success": "created",
		})
	}

}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
