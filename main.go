package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// Task Some columns are not needed in the JSON response
type Task struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	IsCompleted bool      `json:"is_completed"`
	StartTime   time.Time `json:"start_time"`
	Status      string    `json:"status"` // "pending", "started", "paused", "completed"
	TotalTime   int       `json:"total_time"`
}

var db *gorm.DB

func main() {
	dsn := "root:UiiWz5mAdmiygwAbAFem@tcp(117.72.35.68:3306)/timer_me?" +
		"charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	// Migrate the schema to the database automatically
	err1 := db.AutoMigrate(&Task{})
	if err1 != nil {
		return
	}

	r := gin.Default()
	r.Use(cors.Default()) // 添加这行代码以启用默认的CORS中间件

	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.PUT("/tasks/:id/start", startTask)
	r.PUT("/tasks/:id/pause", pauseTask)
	r.PUT("/tasks/:id/complete", completeTask)

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
	// json->struct
	err := c.ShouldBindJSON(&task)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Printf("Error: %s", err.Error())
		return
	}
	log.Printf("Task: %v", task)
	task.StartTime = time.Now()
	// struct->db
	db.Create(&task)

	// 201 Created
	c.JSON(201, task)
}

func updateTask(c *gin.Context) {
	var task Task
	// 因为id是primary key，所以只会有一条记录
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

	// 获取任务ID
	taskID := c.Param("id")
	fmt.Printf("Task ID: %s\n", taskID)

	// 查询数据库
	err := db.Where("id = ?", taskID).First(&task).Error
	if err != nil {
		// 打印错误信息
		fmt.Printf("Error finding task: %s\n", err.Error())
		c.JSON(404, gin.H{
			"error": "record not found",
		})
		return
	}

	// 更新任务状态和开始时间
	task.Status = "started"
	task.StartTime = time.Now()

	// 保存更新到数据库
	db.Save(&task)

	// 返回更新后的任务
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
	task.Duration = (time.Now().Second() - task.StartTime.Second()) + task.Duration
	//task.StartTime = time.Now()
	task.Status = "paused"
	log.Println(task)
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
