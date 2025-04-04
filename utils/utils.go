package utils

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	// "modernc.org/sqlite"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

type CreateSchedule struct {
	Description string `json:"description" binding:"required"`
	Daymonth    string `json:"day_month" binding:"required"`
	Priority    string `json:"prioity" binding:"required"`
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

func InMemory() int16 {
	db, err := gorm.Open(sqlite.Open(os.Getenv("MEMORY")), &gorm.Config{})
	if err != nil {
		return 0
	}
	db.AutoMigrate(&Pin{})

	insertPin := &Pin{TxnPin: 1234}

	db.Create(insertPin)
	return insertPin.TxnPin

	// fmt.Printf("insert ID: %d, Code: %s, Price: %d\n",
	//   insertProduct.ID, insertProduct.Code, insertProduct.Price)
}
