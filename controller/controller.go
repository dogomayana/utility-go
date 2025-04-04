package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	s "strings"
	"sync"
	"time"

	// "strconv"
	// "os"

	"example.com/abobtech/utils"

	// "github.com/joho/godotenv"
	// "github.com/supabase-community/supabase-go"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Sku      string
	ItemName string  `json:"itemName" binding:"required"`
	Unit     int8    `json:"unit" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
	Quantity int8    `json:"quantity" binding:"required"`
}
type client struct {
	count       int
	windowStart time.Time
}
type Schedule struct {
	Description string `json:"description" binding:"required"`
	Daymonth    string `json:"day_month" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
}
type GetUser struct {
	ID          int8   `json:"id"`
	Email       string `json:"email"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Balance     byte   `json:"balance"`
	Device_Type string `json:"device_type"`
}

// FixedWindowRateLimiter returns a Gin middleware that limits requests per client using Fixed Window strategy
func FixedWindowRateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	// Use a thread-safe map to store client data
	clients := make(map[string]*client)
	var mutex sync.Mutex

	// Start a goroutine to clean up old entries
	go func() {
		for {
			time.Sleep(window)
			mutex.Lock()
			for key, c := range clients {
				if time.Since(c.windowStart) > window {
					delete(clients, key)
				}
			}
			mutex.Unlock()
		}
	}()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mutex.Lock()
		cData, exists := clients[clientIP]
		if !exists {
			// First request from this client
			clients[clientIP] = &client{
				count:       1,
				windowStart: time.Now(),
			}
			mutex.Unlock()
			c.Next()
			return
		}

		// Check if the window has passed
		if time.Since(cData.windowStart) > window {
			// Reset the count and window start time
			cData.count = 1
			cData.windowStart = time.Now()
			mutex.Unlock()
			c.Next()
			return
		}

		// Increment the count and check the limit
		if cData.count >= limit {
			// Rate limit exceeded
			mutex.Unlock()
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too Many Requests",
			})
			return
		}

		cData.count++
		mutex.Unlock()
		c.Next()
	}
}

// func ConnectDB() (client *supabase.Client) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("Error loading .env file")
// 	}

// 	supabaseURL := os.Getenv("SUPABASE_URL")
// 	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

// 	client, err = supabase.NewClient(supabaseURL, supabaseAnonKey, nil)

// 	if err != nil {
// 		fmt.Println("cannot initalize client", err)
// 		return
// 	}
// 	return client
// }

func CreateItem(c *gin.Context) {
	// supaClient := ConnectDB()
	var json Item
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "type error or empty value"})
		return
	}
	randN := rand.Intn(10000)
	str := s.Split(json.ItemName, " ")
	sku := fmt.Sprint(randN, str[0])

	fmt.Println(sku)
	row := Item{
		Sku:      sku,
		ItemName: json.ItemName,
		Unit:     json.Unit,
		Price:    json.Price,
		Quantity: json.Quantity,
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": row,
	})
	// _, _, err := supaClient.From("consumables").Insert(row, false, "", "", "").Execute()

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// } else {
	// 	c.JSON(http.StatusCreated, gin.H{
	// 		"success": "created",
	// 	})
	// }
}

// func CreateSchedule(c *gin.Context) {

// 	supaClient := ConnectDB()

// 	var json Schedule
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		c.JSON(http.StatusNotAcceptable, gin.H{"error": "type error or empty value"})
// 		return
// 	}
// 	row := Schedule{
// 		Description: json.Description,
// 		Daymonth:    json.Daymonth,
// 		Priority:    json.Priority,
// 	}
// 	_, _, err := supaClient.From("scheduler").Insert(row, false, "", "", "").Execute()

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	} else {
// 		c.JSON(http.StatusCreated, gin.H{
// 			"success": "created",
// 		})
// 	}

// }
func GetAllSchedules(c *gin.Context) {
	fmt.Printf("%v", utils.InMemory())
	supaClient := utils.DBClient()
	idQuery := c.Query("q")
	var data []byte
	var err error

	if idQuery == "" {
		data, _, err = supaClient.Select("*", "exact", false).Execute()
	} else {
		data, _, err = supaClient.Select("*", "exact", false).Eq("id", idQuery).Execute()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Can't connect to DB",
		})
		return
	}

	var users []GetUser
	err = json.Unmarshal(data, &users)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No Data Found",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
	}
}

func GetSchedule(c *gin.Context) {
	supaClient := utils.DBClient()
	idQuery := c.Query("q")

	data, _, err := supaClient.Select("*", "exact", false).Eq("id", idQuery).Single().Execute()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var schedules GetUser
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

// func DeleteSchedule(c *gin.Context) {

// 	supaClient := ConnectDB()

// 	idQuery := c.Query("q")

// 	_, _, err := supaClient.From("scheduler").Delete("*", "").Eq("id", idQuery).Execute()
// 	if err != nil {

// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"msg": "deleted",
// 		})
// 	}
// }

// func UpdateSchedule(c *gin.Context) {
// 	supaClient := ConnectDB()
// 	var json Schedule
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		c.JSON(http.StatusNotAcceptable, gin.H{"error": "type error or empty value"})

// 		return
// 	}
// 	row := Schedule{
// 		Description: json.Description,
// 		Daymonth:    json.Daymonth,
// 		Priority:    json.Priority,
// 	}

// 	idQuery := c.Query("q")

// 	_, _, err := supaClient.From("scheduler").Update(row, "*", "").Eq("id", idQuery).Execute()

// 	if err != nil {

// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"msg": "Updated",
// 		})
// 	}
// }
