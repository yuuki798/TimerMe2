// main.go
package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type Task struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	IsCompleted bool      `json:"is_completed"`
	StartTime   time.Time `json:"start_time"`
	Status      string    `json:"status"`
	TotalTime   int       `json:"total_time"`
}

var db *gorm.DB

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	initConfig()

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	err1 := db.AutoMigrate(&Task{})
	if err1 != nil {
		return
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.PUT("/tasks/:id/start", startTask)
	r.PUT("/tasks/:id/pause", pauseTask)
	r.PUT("/tasks/:id/complete", completeTask)
	r.PUT("/tasks/:id/reset", resetTask) // 新增reset_time路由

	err3 := r.Run(":8080")
	if err3 != nil {
		return
	}
}

func getTasks(c *gin.Context) {
	var tasks []Task
	db.Find(&tasks)
	c.JSON(200, tasks)
}

func createTask(c *gin.Context) {
	var task Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Printf("Error: %s", err.Error())
		return
	}
	log.Printf("Task: %v", task)
	task.StartTime = time.Now()
	db.Create(&task)
	c.JSON(201, task)
}

func updateTask(c *gin.Context) {
	var task Task
	err := db.Where("id=?", c.Param("id")).First(&task)
	if err.Error != nil {
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}
	db.Save(&task)
	c.JSON(200, task)
}

func deleteTask(c *gin.Context) {
	var task Task
	err := db.Where("id=?", c.Param("id")).First(&task)
	if err.Error != nil {
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}
	db.Delete(&task)
	c.JSON(204, nil)
}

func startTask(c *gin.Context) {
	var task Task
	taskID := c.Param("id")
	fmt.Printf("Task ID: %s\n", taskID)
	err := db.Where("id = ?", taskID).First(&task).Error
	if err != nil {
		fmt.Printf("Error finding task: %s\n", err.Error())
		c.JSON(404, gin.H{
			"error": "record not found",
		})
		return
	}
	if task.Status == "started" {
		c.JSON(400, gin.H{
			"error": "Task already started",
		})
		return
	}
	task.Status = "started"
	task.StartTime = time.Now()
	db.Save(&task)
	c.JSON(200, task)
}

func pauseTask(c *gin.Context) {
	var task Task
	err := db.Where("id=?", c.Param("id")).First(&task).Error
	if err != nil {
		c.JSON(404, gin.H{
			"error": "record not found",
		})
		return
	}
	if task.Status != "started" {
		c.JSON(400, gin.H{
			"error": "Task is not started",
		})
		return
	}
	elapsedTime := int(time.Since(task.StartTime).Seconds())
	task.Duration += elapsedTime
	task.Status = "paused"
	db.Save(&task)
	c.JSON(200, task)
}

func completeTask(c *gin.Context) {
	var task Task
	err := db.Where("id=?", c.Param("id")).First(&task).Error
	if err != nil {
		c.JSON(404, gin.H{
			"error": "record not found",
		})
		return
	}
	task.Status = "completed"
	db.Save(&task)
	c.JSON(200, task)
}

func resetTask(c *gin.Context) { // 新增reset_task函数
	var task Task
	err := db.Where("id=?", c.Param("id")).First(&task).Error
	if err != nil {
		c.JSON(404, gin.H{
			"error": "record not found",
		})
		return
	}
	task.Duration = 0
	task.Status = "pending"
	db.Save(&task)
	c.JSON(200, task)
}
