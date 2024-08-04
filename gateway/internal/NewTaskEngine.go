package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewTaskEngine() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.PUT("/tasks/:id/start", startTask)
	r.PUT("/tasks/:id/pause", pauseTask)
	r.PUT("/tasks/:id/complete", completeTask)
	r.PUT("/tasks/:id/reset", resetTask)
	return r
}
