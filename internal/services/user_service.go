package services

import (
	"context"
	"errors"
	"log"

	"github.com/manabie-com/togo/internal/authorizers"
	"github.com/manabie-com/togo/internal/storages"
)

type UserService interface {
	VerifyUser(ctx context.Context, userID, password string) (string, error)
}

type UserServiceImpl struct {
	Authorizer authorizers.Authorizer
	Store      storages.Store
}

var IncorrectUserIDPassword = errors.New("incorrect user_id/pwd")

func (s *UserServiceImpl) VerifyUser(ctx context.Context, userID, password string) (string, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Could not locate user %s \n", userID)
		return "", IncorrectUserIDPassword
	}
	if password != user.Password {
		log.Printf("Incorrect password in request of user %s \n", userID)
		return "", IncorrectUserIDPassword
	}
	token, err := s.Authorizer.CreateToken(userID)
	if err != nil {
		log.Printf("InternalServerError: could not create token. Error %v \n", err)
		return "", errors.New("Could not create token")
	}

	return token, nil
}
