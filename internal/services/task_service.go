package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/manabie-com/togo/internal/models"
	"github.com/manabie-com/togo/internal/storages"
	"github.com/manabie-com/togo/internal/utils"
)

//TaskService interface
type TaskService interface {
	AddTask(ctx context.Context, task *models.Task, userID string) (*models.Task, error)
	ListTasks(ctx context.Context, userID string, from, to time.Time) (*[]models.Task, error)
}

// TaskServiceImpl implement TaskService interface
type TaskServiceImpl struct {
	Store storages.Store
}

func (s *TaskServiceImpl) AddTask(ctx context.Context, task *models.Task, userID string) (*models.Task, error) {
	now := time.Now().Local()
	task.ID = uuid.New().String()
	task.UserID = userID
	task.CreatedDate = now
	from := utils.BeginningOfDay(now)
	to := utils.BeginningOfDay(now.AddDate(0, 0, 1))

	maxTodo, taskCount, err := s.Store.RetrieveUserMaxTodoAndTaskCount(ctx, userID, from, to)

	if err != nil {
		log.Printf("Could not get user info and tasks. Error: %s \n", err.Error())
		return nil, fmt.Errorf("Could not create task")
	}

	if maxTodo == taskCount {
		log.Printf("User %s already created %d tasks today \n", userID, maxTodo)
		return nil, fmt.Errorf("You are limited to create only %d tasks per day", maxTodo)
	}

	addedTask, err := s.Store.AddTask(ctx, task)
	if err != nil {
		log.Printf("Could not create new task for user %s. Error: %s \n", userID, err.Error())
		return nil, fmt.Errorf("Could not create task")
	}

	return addedTask, nil
}

func (s *TaskServiceImpl) ListTasks(ctx context.Context, userID string, from, to time.Time) (*[]models.Task, error) {
	tasks, err := s.Store.RetrieveTasks(ctx, userID, from, to)
	if err != nil {
		log.Printf("Could not retrive task for user %s. Error: %s \n", userID, err.Error())
		return nil, fmt.Errorf("Could not retrive task for user %s", userID)
	}
	return tasks, nil
}
