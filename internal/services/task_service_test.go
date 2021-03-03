package services

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manabie-com/togo/internal/mocks"
	"github.com/manabie-com/togo/internal/models"
)

func TestAddTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Prepare testing data
	firstUser := "firstUser"
	secondUser := "secondUser"
	task := &models.Task{Content: "Unit test"}

	testCases := []struct {
		name                 string
		userID               string
		expectedErrorMessage error
	}{
		{"Success", firstUser, nil},
		{"Failure - could not add task", firstUser, fmt.Errorf("Could not create task")},
		{"Failure - exceed the max_todo", secondUser, fmt.Errorf("You are limited to create only 1 tasks per day")},
		{"Failure - could not get user info", secondUser, fmt.Errorf("Could not create task")},
	}

	// Create mocked Store and setup expectations
	store := new(mocks.Store)
	store.On("AddTask", ctx, task).Return(task, nil).Once()
	store.On("AddTask", ctx, task).Return(nil, fmt.Errorf("could not add new task")).Once()
	store.On("RetrieveUserMaxTodoAndTaskCount", ctx, firstUser, mock.Anything, mock.Anything).Return(1, 0, nil).Twice()
	store.On("RetrieveUserMaxTodoAndTaskCount", ctx, secondUser, mock.Anything, mock.Anything).Return(1, 1, nil).Once()
	store.On("RetrieveUserMaxTodoAndTaskCount", ctx, secondUser, mock.Anything, mock.Anything).Return(0, 0, fmt.Errorf("could not get user info")).Once()

	taskService := TaskServiceImpl{Store: store}

	for _, ts := range testCases {
		log.Printf("Running unit test add task: %s \n", ts.name)
		addedTask, err := taskService.AddTask(ctx, task, ts.userID)
		if ts.expectedErrorMessage != nil {
			assert.Nil(addedTask)
			assert.Equal(ts.expectedErrorMessage, err)
		} else {
			assert.NotNil(addedTask)
			assert.Nil(err)
		}
	}
}

func TestListTasks(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Prepare testing data
	userID := "firstUser"
	task := models.Task{ID: uuid.New().String(), Content: "Unit test", UserID: userID, CreatedDate: time.Now()}
	tasks := &[]models.Task{task}

	testCases := []struct {
		name                 string
		expectedResponse     *[]models.Task
		expectedErrorMessage error
	}{
		{"Success", tasks, nil},
		{"Failure", nil, fmt.Errorf("Could not retrive task for user %s", userID)},
	}

	// Create mocked Store and setup expectations
	store := new(mocks.Store)
	store.On("RetrieveTasks", ctx, userID, mock.Anything, mock.Anything).Return(tasks, nil).Once()
	store.On("RetrieveTasks", ctx, userID, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("could not get tasks")).Once()

	taskService := TaskServiceImpl{store}

	for _, ts := range testCases {
		log.Printf("Running unit test list tasks: %s \n", ts.name)
		taskList, err := taskService.ListTasks(ctx, userID, time.Now(), time.Now().AddDate(0, 0, 1))
		if ts.expectedResponse != nil {
			assert.Equal(ts.expectedResponse, taskList)
		} else {
			assert.Nil(taskList)
		}
		if ts.expectedErrorMessage != nil {
			assert.Equal(ts.expectedErrorMessage, err)
		} else {
			assert.Nil(err)
		}
	}
}
