package middlewares

import (
	"auth-service/models"
	"auth-service/utils"
	"errors"

	// "fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	DB *gorm.DB
}

func NewAuthMiddleware(DB *gorm.DB) AuthMiddleware {
	return AuthMiddleware{DB}
}

func (am *AuthMiddleware) Authorize(ctx *fiber.Ctx) error {
	cookie_token := ctx.Cookies("token")
	claims := &utils.Claims{}

	token, err := jwt.ParseWithClaims(cookie_token, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fiber.NewError(401, "Your are not authorized")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err == nil && token.Valid {
		var user *models.User
		result := am.DB.Where(&models.User{ID: claims.UserId}).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
				Code:    fiber.StatusNotFound,
				Message: result.Error.Error(),
			})
		}
		ctx.Locals("user", user)
		return ctx.Next()
	}
	return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ServerResponse{
		Code:    fiber.StatusUnauthorized,
		Message: err.Error(),
	})
}
