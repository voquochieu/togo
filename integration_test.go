package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manabie-com/togo/internal/authorizers"
	"github.com/manabie-com/togo/internal/handlers"
	"github.com/manabie-com/togo/internal/mocks"
	"github.com/manabie-com/togo/internal/models"
	"github.com/manabie-com/togo/internal/services"
	"github.com/manabie-com/togo/internal/storages"
)

var token string
var firstUser = "firstUser"
var secondUser = "secondUser"
var correctPass = "example"
var inCorrectPass = "incorrect"

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Prepare testing data
	user := models.User{ID: firstUser, Password: correctPass, MaxTodo: 5}

	testCases := []struct {
		name                 string
		userID               string
		password             string
		expectedCode         int
		expectedErrorMessage string
	}{
		{"Success", firstUser, correctPass, http.StatusOK, ""},
		{"Incorrect password", firstUser, inCorrectPass, http.StatusUnauthorized, `{"error_message":"incorrect user_id/pwd"}`},
		{"Invalid user ID", secondUser, correctPass, http.StatusUnauthorized, `{"error_message":"incorrect user_id/pwd"}`},
	}

	// Create mocked Object and setup expectations
	store := new(mocks.Store)
	store.On("GetUserByID", ctx, firstUser).Return(&user, nil)
	store.On("GetUserByID", ctx, secondUser).Return(nil, fmt.Errorf("Could not locate user %s", secondUser))

	jwtAuthorizer := &authorizers.JWTAuthorizer{
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}
	userService := &services.UserServiceImpl{Authorizer: jwtAuthorizer, Store: store}
	authorizationHandler := &handlers.AuthorizationHandler{Authorizer: jwtAuthorizer, UserService: userService}
	handler := http.Handler(authorizationHandler)

	for _, ts := range testCases {
		log.Printf("Running integration login test case: %s \n", ts.name)
		req, err := http.NewRequest("GET", fmt.Sprintf("/login?user_id=%s&password=%s", ts.userID, ts.password), nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(ts.expectedCode, recorder.Code,
			fmt.Errorf("handler returned wrong status code: got %v, expected %v", recorder.Code, ts.expectedCode))
		if len(ts.expectedErrorMessage) > 0 {
			assert.Equal(ts.expectedErrorMessage, recorder.Body.String(),
				fmt.Errorf("handler returned unexpected error message: got %v, expected %v", recorder.Body.String(), ts.expectedErrorMessage))
		}
		if recorder.Code == http.StatusOK {
			b := recorder.Body.String()
			token = b[10 : len(b)-2]
		}
	}
}

func TestAddTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare testing data
	firstUser := "firstUser"
	task := &models.Task{
		ID:          uuid.New().String(),
		Content:     "Unit test",
		UserID:      firstUser,
		CreatedDate: time.Now(),
	}
	user := &models.User{ID: firstUser, Password: "example", MaxTodo: 5}
	testCases := []struct {
		name                 string
		userID               string
		expectedCode         int
		expectedErrorMessage string
	}{
		{"Success", firstUser, http.StatusOK, ""},
		{"Failure - exceed the max_todo", firstUser, http.StatusTooManyRequests, `{"error_message":"You are limited to create only 1 tasks per day"}`},
	}

	// Create mocked Object and setup expectations
	store := new(mocks.Store)
	store.On("GetUserByID", mock.Anything, firstUser).Return(user, nil)
	store.On("AddTask", mock.Anything, mock.Anything).Return(task, nil).Once()
	store.On("RetrieveUserMaxTodoAndTaskCount", mock.Anything, firstUser, mock.Anything, mock.Anything).Return(1, 0, nil).Once()
	store.On("RetrieveUserMaxTodoAndTaskCount", mock.Anything, firstUser, mock.Anything, mock.Anything).Return(1, 1, nil).Once()

	jwtAuthorizer := &authorizers.JWTAuthorizer{
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}
	userService := &services.UserServiceImpl{Authorizer: jwtAuthorizer, Store: store}
	authorizationHandler := &handlers.AuthorizationHandler{Authorizer: jwtAuthorizer, UserService: userService}
	taskService := &services.TaskServiceImpl{Store: store}
	taskHandler := &handlers.TaskHandler{TaskService: taskService}

	for _, ts := range testCases {
		log.Printf("Running integration add task test case: %s \n", ts.name)
		handler := http.Handler(authorizationHandler.HandleAuthorization(taskHandler))
		req, err := http.NewRequest("POST", "/tasks", io.NopCloser(strings.NewReader(`{"content": "testing"}`)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "applicatiom/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(ts.expectedCode, recorder.Code,
			fmt.Errorf("handler returned wrong status code: got %v, expected %v", recorder.Code, ts.expectedCode))
		if len(ts.expectedErrorMessage) > 0 {
			assert.Equal(ts.expectedErrorMessage, recorder.Body.String(),
				fmt.Errorf("handler returned unexpected body: got %v, expected %v", recorder.Body.String(), ts.expectedErrorMessage))
		}

	}
}

func TestListTasks(t *testing.T) {
	assert := assert.New(t)

	// Prepare testing data
	firstUser := "firstUser"
	now := time.Now()
	createdDate := now.Format("2006-01-02")
	task := &storages.Task{ID: uuid.New().String(), Content: "Unit test", UserID: firstUser, CreatedDate: createdDate}
	tasks := []*storages.Task{task}
	user := models.User{ID: firstUser, Password: "example", MaxTodo: 5}
	testCases := []struct {
		name                 string
		userID               string
		expectedCode         int
		expectedErrorMessage string
	}{
		{"Success", firstUser, http.StatusOK, ""},
	}

	// Create mocked Object and setup expectations
	store := new(mocks.Store)
	store.On("GetUserByID", mock.Anything, firstUser).Return(&user, true)
	store.On("RetrieveTasks", mock.Anything, firstUser, createdDate).Return(tasks, nil).Once()

	jwtAuthorizer := &authorizers.JWTAuthorizer{
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}
	userService := &services.UserServiceImpl{Authorizer: jwtAuthorizer, Store: store}
	authorizationHandler := &handlers.AuthorizationHandler{Authorizer: jwtAuthorizer, UserService: userService}
	taskService := &services.TaskServiceImpl{Store: store}
	taskHandler := &handlers.TaskHandler{TaskService: taskService}

	for _, ts := range testCases {
		log.Printf("Running integration add task test case: %s \n", ts.name)
		handler := http.Handler(authorizationHandler.HandleAuthorization(taskHandler))
		req, err := http.NewRequest("Get", fmt.Sprintf("/tasks?created_date=%s", createdDate), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "applicatiom/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(ts.expectedCode, recorder.Code,
			fmt.Errorf("handler returned wrong status code: got %v, expected %v", recorder.Code, ts.expectedCode))
		if len(ts.expectedErrorMessage) > 0 {
			assert.Equal(ts.expectedErrorMessage, recorder.Body.String(),
				fmt.Errorf("handler returned unexpected body: got %v, expected %v", recorder.Body.String(), ts.expectedErrorMessage))
		}

	}
}
