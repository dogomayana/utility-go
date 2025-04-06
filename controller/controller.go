package controller

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	// "strconv"
	// "os"

	"example.com/abobtech/utils"
	"golang.org/x/crypto/bcrypt"

	// "github.com/joho/godotenv"
	// "github.com/supabase-community/supabase-go"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/mssola/user_agent"
)

type client struct {
	count       int
	windowStart time.Time
}

type GetUser struct {
	ID          int8   `json:"id"`
	Email       string `json:"email"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Balance     byte   `json:"balance"`
	Device_Type string `json:"device_type"`
}
type SignUpT struct {
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Device_Type  string `json:"device_type"`
	Balance      byte   `json:"balance"`
	Auth_Session string `json:"auth_session"`
}
type LogInT struct {
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Auth_Session string `json:"auth_session"`
}
type UpdateAuth struct {
	Auth_Session string `json:"auth_session"`
}

func SignUp(c *gin.Context) {
	userAgent := c.Request.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	device_type := ua.Model()

	supaClient := utils.DBClient()

	var json SignUpT

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type error"})
	}

	byte, err := bcrypt.GenerateFromPassword([]byte(json.Password), 14)
	if err != nil {
	}
	token, _ := utils.JwtTokens(json.Email)
	row := SignUpT{
		FirstName:    json.FirstName,
		LastName:     json.LastName,
		Email:        json.Email,
		Password:     string(byte),
		Device_Type:  device_type,
		Balance:      0.0,
		Auth_Session: token,
	}
	_, _, err = supaClient.Insert(row, false, "", "", "").Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user already exists",
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"success": "created",
		})
	}

}
func LogIn(c *gin.Context) {
	_ = godotenv.Load()
	supaClient := utils.DBClient()

	var jsonBody LogInT
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type error"})
	}

	data, _, err := supaClient.Select("*", "exact", false).Eq("email", jsonBody.Email).Single().Execute()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't connect to DB",
		})
	}
	var users LogInT
	err = json.Unmarshal(data, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't find user",
		})
	}
	var pwd = users.Password
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(jsonBody.Password))

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Wrong email or password"})
	}

	parseValue := utils.ParseJwtToken(users.Auth_Session, users.Email)
	updateAuthSession := UpdateAuth{
		Auth_Session: parseValue,
	}
	if len(parseValue) < 50 {
		c.AbortWithStatusJSON(201, gin.H{
			"status": parseValue,
		})
	} else {
		_, _, err := supaClient.Update(updateAuthSession, "*", "").Eq("email", jsonBody.Email).Execute()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error updating auth session",
			})
		} else {
			c.AbortWithStatusJSON(201, gin.H{
				"status": true,
			})
		}
	}
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

func GetUsers(c *gin.Context) {
	// utils.JwtTokens() // Removed undefined function call
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
