package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	. "net/http"
	_ "test/second"
)

//web socket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *Request) bool {
		return true
	},
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//database
var db *sql.DB
var err error

func conn() {
	driver := "mysql"
	pass := "ayman2371213"
	user := "root"
	name := "godblang"

	db, err = sql.Open(driver, user+":"+pass+"@tcp(localhost:3306)/"+name)

	if err != nil {
		log.Fatal(err)
	}
}

func getAllDB() allEvents {
	row, err := db.Query("SELECT * from godblang.employee")
	if err != nil {
		log.Fatal(err)
	}
	add := event{}
	getAll := allEvents{}
	for row.Next() {
		err := row.Scan(&add.ID, &add.Name, &add.City)
		if err != nil {
			log.Fatal(err)
		}
		getAll = append(getAll, add)
	}
	fmt.Println(getAll)
	return getAll
}
func insertEvent(newEvent event) bool {
	stmt, err := db.Prepare("INSERT INTO  employee (Name,City) values (?,?) ")
	if err != nil {
		log.Fatal(err)
	}
	r, err := stmt.Exec(newEvent.Name, newEvent.City)
	if err != nil {
		log.Fatal(err)
	}
	affectedRows, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The statement affected %d rows\n", affectedRows)
	return true

}

//***********************************//
type event struct {
	ID   string
	Name string
	City string
}

type allEvents []event

var all = allEvents{}

func main() {
	gin.SetMode(gin.ReleaseMode)
	conn()
	r := gin.Default()

	var v1 = r.Group("api/v1")
	{
		v1.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
		v1.GET("/events", getAllEvents)
		v1.POST("/addEvent", addEvent)
	}

	r.Run()
}

func addEvent(c *gin.Context) {
	var newEvent event
	newEvent.ID = c.PostForm("ID")
	newEvent.City = c.PostForm("City")
	newEvent.Name = c.PostForm("Name")
	var inserted bool
	inserted = insertEvent(newEvent)
	if inserted {
		fmt.Println("event:: ", newEvent)
		c.JSON(201, gin.H{
			"sucess":  true,
			"message": "event has created successfully",
			"event":   newEvent,
		})
	} else {
		fmt.Println("event:: ", newEvent)
		c.JSON(201, gin.H{
			"sucess":  false,
			"message": "not created",
			"event":   newEvent,
		})
	}
}
func getAllEvents(c *gin.Context) {
	all = getAllDB()
	fmt.Println("all Events : ", all)
	c.JSON(200, gin.H{
		"data": all,
	})
}
