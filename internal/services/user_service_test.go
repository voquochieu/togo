package services

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/manabie-com/togo/internal/mocks"
	"github.com/manabie-com/togo/internal/models"
)

func TestVerifyUser(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Prepare testing data
	firstUser := "firstUser"
	secondUser := "secondUser"
	password := "example"
	user := models.User{ID: firstUser, Password: password, MaxTodo: 5}

	testCases := []struct {
		name                 string
		userID               string
		password             string
		expectedErrorMessage error
	}{
		{"Success", firstUser, password, nil},
		{"Failure - could not generate token", firstUser, password, fmt.Errorf("Could not create token")},
		{"Failure - incorrect password", firstUser, "wrong_pass", fmt.Errorf("incorrect user_id/pwd")},
		{"Failure - invalid userID", secondUser, password, fmt.Errorf("incorrect user_id/pwd")},
	}

	// Create mocked Store and setup expectations
	store := new(mocks.Store)
	store.On("GetUserByID", ctx, firstUser).Return(&user, nil)
	store.On("GetUserByID", ctx, secondUser).Return(nil, fmt.Errorf("Could not locate user %s", secondUser))

	authorizer := new(mocks.Authorizer)
	authorizer.On("CreateToken", firstUser).Return("generated_token", nil).Once()
	authorizer.On("CreateToken", firstUser).Return("", fmt.Errorf("could not generate token")).Once()

	userService := UserServiceImpl{authorizer, store}

	for _, ts := range testCases {
		log.Printf("Running unit test verify user: %s \n", ts.name)
		token, err := userService.VerifyUser(ctx, ts.userID, ts.password)
		if ts.expectedErrorMessage != nil {
			assert.Empty(token)
			assert.Equal(ts.expectedErrorMessage, err)
		} else {
			assert.NotNil(token)
			assert.Nil(err)
		}
	}
}
