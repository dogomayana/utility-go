package utils

import (
	// "database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO

	// "modernc.org/sqlite"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

type MyCustomClaims struct {
	Auth_Name string `json:"auth_name"`
	jwt.RegisteredClaims
}

type Pin struct {
	TxnPin int16 `json:"txn_pin" binding:"required"`
}

func DBClient() *postgrest.QueryBuilder {
	_ = godotenv.Load()

	url := os.Getenv("DB_URL")
	key := os.Getenv("DB_ANON_KEY")
	table := os.Getenv("DB_TABLE")
	client, err := supabase.NewClient(url, key, nil)

	if err != nil {
		fmt.Println("cannot initalize client", err)
		return nil
	}
	return client.From(table)
}

func InMemory() {
	db, err := gorm.Open(sqlite.Open(os.Getenv("MEMORY")), &gorm.Config{})
	// db, err := sql.Open("sqlite", os.Getenv("MEMORY"))
	if err != nil {
		return
	}

	db.AutoMigrate(&Pin{})

	insertPin := &Pin{TxnPin: 1234}

	db.Create(insertPin)
	// return insertPin.TxnPin

	fmt.Printf("insert ID: %d", insertPin.TxnPin)

}

func JwtTokens(email string) (string, error) {
	_ = godotenv.Load()

	mySigningKey := []byte(os.Getenv("SIGNIN_TOKEN"))
	claims := MyCustomClaims{
		"authSessionKey",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func ParseJwtToken(tokenString, email string) string {
	_ = godotenv.Load()

	mySigningKey := []byte(os.Getenv("SIGNIN_TOKEN"))

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		newToken, _ := JwtTokens(email)
		return newToken
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		return claims.Issuer
	} else {
		return ""
	}
}

// supaClient := ConnectDB()
// now := time.Now()
// seconds := now.Second()
// day := now.Day()
// rn := rand.Intn(99-10+1) + 10
// monthStr := now.Format("01")
// dayStr := fmt.Sprintf("%02d", day)
// secondsStr := fmt.Sprintf("%02d", seconds)
// ramdomStr := strconv.Itoa(rn)
