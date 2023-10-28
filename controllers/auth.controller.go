package controllers

import (
	"auth-service/models"
	"auth-service/utils"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthCotroller(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

/***** Register a user ******/
func (ac *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	var payload *utils.RegisterUserInput

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ServerResponse{
			Success: false,
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	newUser := models.User{
		Email:          payload.Email,
		Username:       payload.Username,
		Password:       payload.Password,
		PolicyId:       payload.PolicyId,
		OrganizationId: payload.OrganizationId,
		Profile: models.Profile{
			Bio: "update your bio",
		},
	}

	result := ac.DB.Create(&newUser)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return ctx.Status(fiber.StatusConflict).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusConflict,
		})
	} else if result.Error != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.ServerResponse{
		Success: true,
		Message: "Registered successfully",
		Code:    fiber.StatusCreated,
	})
}

/***** Get users ******/
func (ac *AuthController) GetUsers(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Locals("user"))
	var users []utils.UserOutput
	result := ac.DB.Model(&models.User{}).Find(&users)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "OK",
		Code:    fiber.StatusOK,
		Data:    users,
	})
}

/****** Get a user by id ******/
func (ac *AuthController) GetUserById(ctx *fiber.Ctx) error {
	var user *models.User
	userId := uuid.FromStringOrNil(ctx.Params("userId"))
	result := ac.DB.Where(&models.User{ID: userId}).Preload("Profile").First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	} else if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "OK",
		Code:    fiber.StatusOK,
		Data:    new(utils.UserOutput).FromUser(user),
	})
}

/****** Log user in ******/
func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var payload *utils.LoginUserInput

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ServerResponse{
			Success: false,
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	var user *models.User
	result := ac.DB.Where(
		&models.User{
			Username: payload.Username,
			Password: payload.Password,
		},
	).Preload("Profile").First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	} else if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	claims := &utils.Claims{
		UserId:         user.ID,
		OrganizationId: user.OrganizationId,
		PolicyId:       user.PolicyId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    "devsware.net",
			Subject:   "berny@mail.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := os.Getenv("SECRET_KEY")
	str, err := token.SignedString([]byte(key))

	if err != nil {
		return ctx.Status(fiber.StatusPreconditionFailed).JSON(utils.ServerResponse{
			Success: false,
			Message: err.Error(),
			Code:    fiber.StatusPreconditionFailed,
		})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = str

	ctx.Cookie(cookie)

	if user.Profile.TwoFactorEnabled {
		return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
			Success: true,
			Message: "user logged in successfully",
			Code:    fiber.StatusOK,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "OK",
		Code:    fiber.StatusOK,
		Data:    new(utils.UserOutput).FromUser(user),
	})
}

/****** Verify totp ******/
func (ac *AuthController) VerifyTotp(ctx *fiber.Ctx) error {
	//needs a valid token, userid, topt-token
	var payload *utils.TOTPInput

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ServerResponse{
			Success: false,
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	var user *models.User
	result := ac.DB.Where(&models.User{ID: payload.UserId}).First(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	valid := totp.Validate(payload.Token, user.TotpSecret)

	if !valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ServerResponse{
			Success: false,
			Message: "invalid code",
			Code:    fiber.StatusUnauthorized,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "OK",
		Code:    fiber.StatusOK,
		Data:    new(utils.UserOutput).FromUser(user),
	})
}

/***** Enable totp ******/
func (ac *AuthController) EnableTotp(ctx *fiber.Ctx) error {
	// verify api token with returns a user
	userId := uuid.FromStringOrNil(ctx.Params("userId"))
	var user *models.User

	result := ac.DB.Where(&models.User{ID: userId}).Preload("Profile").First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	// 2fa already enabled
	if user.Profile.TwoFactorEnabled {
		return ctx.Status(fiber.StatusOK).JSON(&utils.TOTPOut{
			UserId:     user.ID,
			TotpSecret: user.TotpSecret,
			TotpUrl:    user.TotpUrl,
		})
	}

	key, otp_err := totp.Generate(totp.GenerateOpts{
		Issuer:      "devsware.net",
		AccountName: "berny@mail.com",
		SecretSize:  15,
	})

	if otp_err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ServerResponse{
			Success: false,
			Message: otp_err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	user.TotpSecret = key.Secret()
	user.TotpUrl = key.URL()
	user.Profile.TwoFactorEnabled = true

	ac.DB.Save(&user)

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "OK",
		Code:    fiber.StatusOK,
		Data: &utils.TOTPOut{
			UserId:     user.ID,
			TotpSecret: user.TotpSecret,
			TotpUrl:    user.TotpUrl,
		},
	})
}

/****** Disable totp ******/
func (ac *AuthController) DisableTotp(ctx *fiber.Ctx) error {
	// verify api token with returns a user
	userId := uuid.FromStringOrNil(ctx.Params("userId"))
	var profile *models.Profile
	result := ac.DB.Where(&models.Profile{UserId: userId}).First(&profile)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	profile.TwoFactorEnabled = false
	ac.DB.Save(profile)

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "2fa disabled",
		Code:    fiber.StatusOK,
	})
}

/****** Verify email ******/
func (ac *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	// verify api token with returns a user
	userId := uuid.FromStringOrNil(ctx.Params("userId"))
	var profile *models.Profile
	result := ac.DB.Where(&models.Profile{UserId: userId}).First(&profile)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ServerResponse{
			Success: false,
			Message: result.Error.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	profile.EmailVerified = true
	ac.DB.Save(profile)

	return ctx.Status(fiber.StatusOK).JSON(utils.ServerResponse{
		Success: true,
		Message: "your email is verified",
		Code:    fiber.StatusOK,
	})
}
