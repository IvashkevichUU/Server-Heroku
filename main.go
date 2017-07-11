package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	_ "github.com/lib/pq"
)

var (
	repeat int
	db     *sql.DB
)

func repeatFunc(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go блеа!")
	}
	c.String(http.StatusOK, buffer.String())
}

func createdb(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS students (id SERIAL NOT NULL, fio CHARACTER VARYING(300) NOT NULL, info TEXT NOT NULL, score INTEGER NOT NULL )"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}
	if _, err := db.Exec("INSERT INTO students (fio, info, score) VALUES ('Vasily Romanov', 'company: Mail.ru Group', '10')"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}
}

// PrintByID print student by id
func PrintByID(c *gin.Context) {
	rows, err := db.Query("SELECT fio FROM students")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading students: %q", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick string
		if err := rows.Scan(&tick); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning students: %q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick))
	}
}

func dbFunc(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning ticks: %q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error
	tStr := os.Getenv("REPEAT")
	repeat, err = strconv.Atoi(tStr)
	if err != nil {
		log.Print("Error converting $REPEAT to an int: %q - Using default", err)
		repeat = 5
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatFunc)
	router.GET("/db", dbFunc)

	router.GET("/students", createdb)

	router.GET("/studentid", PrintByID)

	router.Run(":" + port)
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}




//package main
//
//import (
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/gin-gonic/gin"
//	"github.com/russross/blackfriday"
//)
//
//func main() {
//	port := os.Getenv("PORT")
//
//	if port == "" {
//		log.Fatal("$PORT must be set")
//	}
//
//	router := gin.New()
//	router.Use(gin.Logger())
//	router.LoadHTMLGlob("templates/*.tmpl.html")
//	router.Static("/static", "static")
//
//	router.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "index.tmpl.html", nil)
//	})
//
//	router.GET("/mark", func(c *gin.Context) {
//		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
//	})
//
//	router.Run(":" + port)
//}
//package main
//
//import (
//	"bytes"
//	"log"
//	"net/http"
//	"os"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	"github.com/russross/blackfriday"
//)
//
//var (
//	repeat int
//)
//
//func repeatHandler(c *gin.Context) {
//	var buffer bytes.Buffer
//	for i := 0; i < repeat; i++ {
//		buffer.WriteString("Hello from Go!\n")
//	}
//	c.String(http.StatusOK, buffer.String())
//}
//
//func main() {
//	var err error
//	port := os.Getenv("PORT")
//
//	if port == "" {
//		log.Fatal("$PORT must be set")
//	}
//
//	tStr := os.Getenv("REPEAT")
//	repeat, err = strconv.Atoi(tStr)
//	if err != nil {
//		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
//		repeat = 5
//	}
//
//	router := gin.New()
//	router.Use(gin.Logger())
//	router.LoadHTMLGlob("templates/*.tmpl.html")
//	router.Static("/static", "static")
//
//	router.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "index.tmpl.html", nil)
//	})
//
//	router.GET("/mark", func(c *gin.Context) {
//		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
//	})
//
//	router.GET("/repeat", repeatHandler)
//
//	router.Run(":" + port)
//}