package routes

import (
	"github.com/nturu/microservice-template/internal/handlers"
	"github.com/nturu/microservice-template/internal/middleware"
	"github.com/nturu/microservice-template/internal/repository"
	"github.com/nturu/microservice-template/internal/validators"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func registerAuth(router fiber.Router, db *gorm.DB) {
	authRouter := router.Group("auth")
	userRepo := repository.NewUserRepository(db)
	handler := handlers.NewAuthHandler(userRepo)

	authRouter.Post("/signup", validators.ValidateRegisterUserSchema, handler.Register)
	authRouter.Post("/signin", validators.ValidateLoginUser, handler.Authenticate)

	authRouter.Get("/password-reset/send-otp", handler.SendOTPForPasswordReset)
	authRouter.Post("/verify-otp/:email/:otp", handler.VerifyOTPAndGenerateToken)
	authRouter.Post("/password-reset/new-password", validators.ValidatePasswordReset, middleware.JWTMiddleware(db), handler.ResetPassword)
}
