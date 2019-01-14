package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Tasks struct {
	Tid      int    `gorm:"primary_key" form:"tid" json:"tid"`
	Taskname string `gorm:"not null" form:"taskname" json:"taskname"`
	Status   bool   `gorm:"not null" form:"status" json:"status"`
}

func main() {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.POST("/tasks", PostUser)
		v1.GET("/tasks", GetAllTasks)
		v1.GET("/tasks/:tid", GetTask)
		v1.POST("/tasks/:tid", MarkTask)
	}

	r.Run(":8080")
}
func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	db.LogMode(true)
	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	// if !db.HasTable(&Users{}) {
	//     db.CreateTable(&Users{})
	//     db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Users{})
	// }
	if !db.HasTable(&Tasks{}) {
		db.CreateTable(&Tasks{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tasks{})
	}

	return db
}
func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var task Tasks
	c.Bind(&task)

	// fmt.Println(task)
	if task.Taskname != "" {
		db.Create(&task)
		c.JSON(201, gin.H{"success": task})
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

}
func GetAllTasks(c *gin.Context) {
	db := InitDb()
	defer db.Close()
	var task []Tasks
	db.Find(&task)
	c.JSON(200, task)
}

func GetTask(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	id := c.Params.ByName("tid")
	temp, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		//
	}
	var task Tasks
	db.First(&task, temp)
	fmt.Println(task)
	if task.Tid != 0 {
		c.JSON(200, task)
	} else {
		c.JSON(404, gin.H{"error": "Task not found"})
	}
}
func MarkTask(c *gin.Context) {
	db := InitDb()
	defer db.Close()
	id := c.Params.ByName("tid")
	temp, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		//
	}
	var task Tasks
	db.First(&task, temp)

	if task.Tid != 0 {
		db.Delete(&task)
		c.JSON(200, gin.H{"success": "Task #" + id + " completed and deleted"})
	} else {
		c.JSON(404, gin.H{"error": "Task not found"})
	}
}
