package database

// 	// "database/sql"
// "fmt"

// 	// "os"
// 	// "strconv"

// 	// "github.com/joho/godotenv"
// 	// _ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
// 	"github.com/supabase-community/supabase-go"

func connectDB() {
	// client, err := supabase.NewClient("https://ciscchstkoanleiqhyiu.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImNpc2NjaHN0a29hbmxlaXFoeWl1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2ODA3MzU1NzgsImV4cCI6MTk5NjMxMTU3OH0.Za39xLn4BQv7U_2Of_IUYmuv_x8rLke19GDc52TkDv4", nil)

	// if err != nil {
	// 	fmt.Println("cannot initalize client", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "cannot initalize client",
	// 	})
	// 	return
	// }
	// return client
}

// // var Version = "1.0"

// // var Db *sql.DB

// // func ConnectDatabase() {

// 	// err := godotenv.Load() //by default, it is .env so we don't have to write
// 	// if err != nil {
// 	// 	fmt.Println("Error is occurred  on .env file please check")
// 	// }
// 	//we read our .env file
// 	// host := os.Getenv("HOST")
// 	// port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
// 	// user := os.Getenv("USER")
// 	// dbname := os.Getenv("DB_NAME")
// 	// pass := os.Getenv("PASSWORD")

// 	// // set up postgres sql to open it.
// 	// psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
// 	// 	host, port, user, dbname, pass)
// 	// db, errSql := sql.Open("postgres", psqlSetup)
// 	// if errSql != nil {
// 	// 	fmt.Println("There is an error while connecting to the database ", err)
// 	// 	panic(err)
// 	// } else {
// 	// 	Db = db
// 	// 	fmt.Println("Successfully connected to database!")
// 	// }
// // }
//
