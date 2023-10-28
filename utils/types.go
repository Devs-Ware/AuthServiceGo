package utils

import (
	"auth-service/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type RegisterUserInput struct {
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Username       string    `json:"username"`
	PolicyId       uuid.UUID `json:"policyId"`
	OrganizationId uuid.UUID `json:"organizationId"`
}

type UserOutput struct {
	ID             uuid.UUID      `json:"id"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	PolicyId       uuid.UUID      `json:"policyId"`
	OrganizationId uuid.UUID      `json:"organizationId"`
	ProfileId      uuid.UUID      `json:"profileId"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt"`
}

func (u *UserOutput) FromUser(user *models.User) *UserOutput {
	u.ID = user.ID
	u.Email = user.Email
	u.Username = user.Username
	u.PolicyId = user.PolicyId
	u.OrganizationId = user.OrganizationId
	u.ProfileId = user.ProfileId
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt
	return u
}

type LoginUserInput struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type TOTPInput struct {
	UserId uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
}

type TOTPOut struct {
	UserId     uuid.UUID `json:"userId"`
	TotpSecret string    `json:"totpSecret"`
	TotpUrl    string    `json:"totpUrl"`
}

type ServerResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
}

type Claims struct {
	UserId         uuid.UUID `json:"userId"`
	OrganizationId uuid.UUID `json:"organizationId"`
	PolicyId       uuid.UUID `json:"policyId"`
	jwt.RegisteredClaims
}
