package utils

import (
	"auth-service/models"
	"context"
	"fmt"
	"time"
)

type LoggingService struct {
	next Service
}

func NewLoggingService(next Service) Service {
	return &LoggingService{
		next: next,
	}
}

func (s *LoggingService) GetUser(ctx context.Context) (user *models.User, err error) {
	defer func(start time.Time) {
		fmt.Printf("user=%v err=%s took=%v\n", user, err, time.Since(start))
	}(time.Now())

	return s.next.GetUser(ctx)
}
