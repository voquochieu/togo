package storages

import (
	"context"
	"time"

	"github.com/manabie-com/togo/internal/models"
)

// Store DB interface
type Store interface {
	RetrieveTasks(ctx context.Context, userID string, from, to time.Time) (*[]models.Task, error)
	AddTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	RetrieveUserMaxTodoAndTaskCount(ctx context.Context, userID string, from, to time.Time) (int, int, error)
}
