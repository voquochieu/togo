package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/manabie-com/togo/internal/models"
)

var users = []models.User{
	{
		ID:       "firstUser",
		Password: "example",
	},
	{
		ID:       "secondUser",
		Password: "example",
	},
}

var tasks = []models.Task{
	{
		Content: "Task 1",
	},
	{
		Content: "Task 2",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Task{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Task{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Task{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].ID = uuid.New().String()
		tasks[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
