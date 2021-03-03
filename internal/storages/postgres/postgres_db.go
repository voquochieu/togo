package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/manabie-com/togo/internal/models"
)

type (
	PostgresDB struct {
		DB *gorm.DB
	}
	UserLimitAndTasks struct {
		MaxTodo   int
		TaskCount int
	}
)

func (d *PostgresDB) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user := &models.User{}
	err := d.DB.Debug().Model(models.User{}).Where("id = ?", userID).Take(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return &models.User{}, fmt.Errorf("User Not Found")
	}
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (d *PostgresDB) RetrieveTasks(ctx context.Context, userID string, from, to time.Time) (*[]models.Task, error) {
	tasks := []models.Task{}
	err := d.DB.Debug().Model(models.Task{}).Where("user_id = ? and created_date >= ? and created_date <= ?", userID, from, to).Find(&tasks).Error
	if gorm.IsRecordNotFoundError(err) {
		return &[]models.Task{}, nil
	}
	if err != nil {
		return &[]models.Task{}, err
	}
	return &tasks, nil
}

func (d *PostgresDB) AddTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	err := d.DB.Debug().Create(&task).Error
	if err != nil {
		return &models.Task{}, err
	}
	return task, nil
}

func (d *PostgresDB) RetrieveUserMaxTodoAndTaskCount(ctx context.Context, userID string, from, to time.Time) (int, int, error) {
	result := &UserLimitAndTasks{}
	err := d.DB.Debug().Table("users").Select("users.max_todo, COUNT(tasks.id) AS task_count").
		Joins("LEFT JOIN tasks ON tasks.user_id = users.id").
		Where("users.id = ? and created_date >= ? and created_date < ?", userID, from, to).
		Group("users.max_todo").
		Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return -1, -1, fmt.Errorf("User Not Found")
	}
	if err != nil {
		return -1, -1, err
	}
	return result.MaxTodo, result.TaskCount, nil
}
