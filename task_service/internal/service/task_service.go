package service

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"task_service/internal/entity"
	"time"

	"gorm.io/gorm"
	pb "proto/task"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
func InitDB() {
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
		panic("failed to connect database")
	}
	db.AutoMigrate(&entity.Task{})
}

type Server struct {
	pb.UnimplementedTaskServiceServer
}

var db *gorm.DB

func (s *Server) GetTasks(ctx context.Context, in *pb.Empty) (*pb.TaskList, error) {
	var tasks []entity.Task
	db.Find(&tasks)
	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, &pb.Task{
			Id:          uint64(task.ID),
			Name:        task.Name,
			Duration:    int32(task.Duration),
			IsCompleted: task.IsCompleted,
			StartTime:   task.StartTime.Format(time.DateTime),
			Status:      task.Status,
			TotalTime:   int32(task.TotalTime),
		})
	}
	return &pb.TaskList{Tasks: pbTasks}, nil
}

func (s *Server) CreateTask(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	task := entity.Task{
		Name:        in.Name,
		Duration:    int(in.Duration),
		IsCompleted: in.IsCompleted,
		Status:      in.Status,
		StartTime:   time.Now(),
		TotalTime:   int(in.TotalTime),
	}
	db.Create(&task)
	in.Id = uint64(task.ID)
	return in, nil
}

func (s *Server) UpdateTask(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	task.Name = in.Name
	task.Duration = int(in.Duration)
	task.IsCompleted = in.IsCompleted
	task.Status = in.Status
	task.TotalTime = int(in.TotalTime)
	db.Save(&task)
	return in, nil
}

func (s *Server) DeleteTask(ctx context.Context, in *pb.TaskId) (*pb.Empty, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	db.Delete(&task)
	return &pb.Empty{}, nil
}

func (s *Server) StartTask(ctx context.Context, in *pb.TaskId) (*pb.Task, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	if task.Status == "started" {
		return nil, fmt.Errorf("task already started")
	}
	task.Status = "started"
	task.StartTime = time.Now()
	db.Save(&task)
	return &pb.Task{
		Id:          uint64(task.ID),
		Name:        task.Name,
		Duration:    int32(task.Duration),
		IsCompleted: task.IsCompleted,
		StartTime:   task.StartTime.Format(time.RFC3339),
		Status:      task.Status,
		TotalTime:   int32(task.TotalTime),
	}, nil
}

func (s *Server) PauseTask(ctx context.Context, in *pb.TaskId) (*pb.Task, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	if task.Status != "started" {
		return nil, fmt.Errorf("task is not started")
	}
	elapsedTime := int(time.Since(task.StartTime).Seconds())
	task.Duration += elapsedTime
	task.Status = "paused"
	db.Save(&task)
	return &pb.Task{
		Id:          uint64(task.ID),
		Name:        task.Name,
		Duration:    int32(task.Duration),
		IsCompleted: task.IsCompleted,
		StartTime:   task.StartTime.Format(time.RFC3339),
		Status:      task.Status,
		TotalTime:   int32(task.TotalTime),
	}, nil
}

func (s *Server) CompleteTask(ctx context.Context, in *pb.TaskId) (*pb.Task, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	task.Status = "completed"
	db.Save(&task)
	return &pb.Task{
		Id:          uint64(task.ID),
		Name:        task.Name,
		Duration:    int32(task.Duration),
		IsCompleted: task.IsCompleted,
		StartTime:   task.StartTime.Format(time.RFC3339),
		Status:      task.Status,
		TotalTime:   int32(task.TotalTime),
	}, nil
}

func (s *Server) ResetTask(ctx context.Context, in *pb.TaskId) (*pb.Task, error) {
	var task entity.Task
	err := db.First(&task, in.Id).Error
	if err != nil {
		return nil, err
	}
	task.Duration = 0
	task.Status = "pending"
	task.IsCompleted = false
	task.StartTime = time.Now() // 这里设置为当前时间
	db.Save(&task)
	fmt.Println(task)
	pbTask := &pb.Task{
		Id:          uint64(task.ID),
		Name:        task.Name,
		Duration:    int32(task.Duration),
		IsCompleted: task.IsCompleted,
		StartTime:   task.StartTime.Format(time.RFC3339),
		Status:      task.Status,
		TotalTime:   int32(task.TotalTime),
	}
	fmt.Printf("Reset Task: %+v\n", pbTask) // 打印重置后的任务
	return pbTask, nil
}
