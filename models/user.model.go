package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uuid.UUID `json:"id" gorm:"type=uuid;primary_key;"`
	Email          string    `json:"email" gorm:"uniqueIndex;not null"`
	Password       string    `json:"password" gorm:"type=varchar(255);not null"`
	Username       string    `json:"username" gorm:"type=varchar(255);not null"`
	PolicyId       uuid.UUID `json:"policyId" gorm:"type=uuid"`
	OrganizationId uuid.UUID `json:"organizationId" gorm:"type=uuid"`
	TotpSecret     string    `json:"totpSecret" gorm:"type=varchar(255);not null"`
	TotpUrl        string    `json:"totpUrl" gorm:"type=varchar(255);not null"`
	ProfileId      uuid.UUID `json:"profileId"`
	Profile        Profile   `json:"profile" gorm:"foreignKey:UserId"`
}

type Profile struct {
	gorm.Model
	ID               uuid.UUID `json:"id" gorm:"type=uuid;primary_key"`
	EmailVerified    bool      `json:"emailVerified" gorm:"default:false"`
	TwoFactorEnabled bool      `json:"twoFactorEnabled" gorm:"default:false"`
	CreatedAt        time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	Bio              string    `json:"bio"`
	UserId           uuid.UUID `json:"userId" gorm:"type=uuid"`
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.NewV4()
	user.Profile.ID = uuid.NewV4()
	user.ProfileId = user.Profile.ID
	return nil
}

