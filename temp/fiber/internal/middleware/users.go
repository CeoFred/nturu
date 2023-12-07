package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/nturu/microservice-template/constants"
	"gorm.io/gorm"

	"github.com/nturu/microservice-template/internal/helpers"
	"github.com/nturu/microservice-template/internal/models"
	"github.com/nturu/microservice-template/internal/repository"

	"strings"
)

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
func NewError(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

var (
	constant = constants.New()
)

func OnlyAdmin(db *gorm.DB, u *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		claims, _ := c.Locals("claims").(*helpers.AuthTokenJwtClaim)

		user, _, err := u.FindUserByCondition("user_id", claims.UserId)
		if err != nil {
			return helpers.Dispatch500Error(c, err)
		}
		if user.Role != models.AdminRole {
			return helpers.Dispatch500Error(c, NewError("Unauthorized access"))
		}
		return c.Next()
	}
}

func JWTMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Extract the token from the "Bearer <jwt>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return fiber.ErrUnauthorized
		}

		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			// Provide the same JWT secret key used for signing the tokens
			return []byte(constant.JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			return fiber.ErrUnauthorized
		}

		// Attach the claims to the request context for further use
		c.Locals("claims", claims)

		_, found, err := repository.NewUserRepository(db).FindUserByCondition("email", claims.Email)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		if !found {
			return fiber.ErrUnauthorized
		}

		// Proceed to the next middleware or route handler
		return c.Next()
	}
}
