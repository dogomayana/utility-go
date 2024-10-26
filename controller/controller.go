package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strconv"

	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"

	"github.com/gin-gonic/gin"
)

type Schedule struct {
	Description string `json:"description" binding:"required"`
	Daymonth    string `json:"day_month" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
}
type GetSchedules struct {
	ID          int8   `json:"id"`
	Description string `json:"description"  `
	Daymonth    string `json:"day_month" `
	Priority    string `json:"priority"  `
}

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

func CreateSchedule(c *gin.Context) {

	supaClient := ConnectDB()

	var json Schedule
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "type error"})
		return
	}
	row := Schedule{
		Description: json.Description,
		Daymonth:    json.Daymonth,
		Priority:    json.Priority,
	}
	_, _, err := supaClient.From("scheduler").Insert(row, false, "", "", "").Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"success": "created",
		})
	}

}
func GetAllSchedules(c *gin.Context) {
	supaClient := ConnectDB()
	data, _, err := supaClient.From("scheduler").Select("*", "exact", false).Execute()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	var schedules []GetSchedules
	err = json.Unmarshal(data, &schedules)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "server error",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": schedules,
		})
	}
}

func DeleteSchedule(c *gin.Context) {

	supaClient := ConnectDB()

	idQuery := c.Query("q")

	_, _, err := supaClient.From("scheduler").Delete("*", "").Eq("id", idQuery).Execute()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "deleted",
		})
	}
}

func UpdateSchedule(c *gin.Context) {
	supaClient := ConnectDB()
	var json Schedule
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "type error"})
		return
	}
	row := Schedule{
		Description: json.Description,
		Daymonth:    json.Daymonth,
		Priority:    json.Priority,
	}

	idQuery := c.Query("q")

	_, _, err := supaClient.From("scheduler").Update(row, "*", "").Eq("id", idQuery).Execute()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Updated",
		})
	}
}
