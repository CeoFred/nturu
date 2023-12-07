package routes

import (
	"github.com/nturu/microservice-template/internal/handlers"
	"github.com/nturu/microservice-template/internal/middleware"
	"github.com/nturu/microservice-template/internal/repository"
	"github.com/nturu/microservice-template/internal/validators"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func registerUser(router fiber.Router, db *gorm.DB) {
	userRouter := router.Group("users")
	handler := handlers.NewUserHandler(repository.NewUserRepository(db))

	userRouter.Get("/profile", middleware.JWTMiddleware(db), handler.UserProfile)
	userRouter.Post("/logo", middleware.JWTMiddleware(db), handler.UploadLogo)
	userRouter.Post("/file-upload", middleware.JWTMiddleware(db), handler.FileUpload)
	userRouter.Put("/", middleware.JWTMiddleware(db), validators.ValidateUpdateUserProfile, handler.UpdateUserProfile)
}
