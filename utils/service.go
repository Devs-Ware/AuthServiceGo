package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"auth-service/models"
	"net/http"
)

type Service interface {
	GetUser(context.Context) (*models.User, error)
}

type UserService struct {
	url string
}

func NewUserService(url string) Service {
	return &UserService{
		url: url,
	}
}

func (s *UserService) GetUser(ctx context.Context) (*models.User, error) {
	res, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	user := &models.User{}
	if err := json.NewDecoder(res.Body).Decode(user); err != nil {
		fmt.Println("error here")
		return nil, err
	}

	fmt.Printf("user: %v\n", user.Username)

	return user, nil
}
