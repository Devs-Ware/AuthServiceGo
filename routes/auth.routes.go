package routes

import (
	ctl "auth-service/controllers"
	// "auth-service/middlewares"
	mw "auth-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

type RouteController struct {
	authController ctl.AuthController
	authMiddleware mw.AuthMiddleware
}

func NewRouteController(authController ctl.AuthController, authMiddleware mw.AuthMiddleware) RouteController {
	return RouteController{authController, authMiddleware}
}

func (rc *RouteController) AuthRoute(route fiber.Router) {
	auth := route.Group("auth")

	auth.Post("/register", rc.authController.RegisterUser)
	auth.Post("/login", rc.authController.Login)
	auth.Get("/users", rc.authMiddleware.Authorize, rc.authController.GetUsers)
	auth.Get("/users/:userId<guid>", rc.authController.GetUserById)
	auth.Patch("/users/:userId<guid>/totp/verify", rc.authController.VerifyTotp)
	auth.Patch("/users/:userId<guid>/totp/enable", rc.authController.EnableTotp)
	auth.Patch("/users/:userId<guid>/totp/disable", rc.authController.DisableTotp)
	auth.Patch("/users/:userId<guid>/email/verify", rc.authController.VerifyEmail)
}
