package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nturu/microservice-template/constants"
	"github.com/nturu/microservice-template/internal/helpers"
	"github.com/nturu/microservice-template/internal/models"
	"github.com/nturu/microservice-template/internal/otp"
	"github.com/nturu/microservice-template/internal/repository"
	"github.com/nturu/microservice-template/sendgrid"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type AuthHandler struct {
	userRepository *repository.UserRepository
}

func NewAuthHandler(
	userRepo *repository.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		userRepository: userRepo,
	}
}

var (
	constant = constants.New()
)

type RegisterResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    RegisterResponseData `json:"data"`
}

type RegisterResponseData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
type InputCreateUser struct {
	BuisnessName string             `json:"business_name" validate:"required"`
	Manager      string             `json:"manager" validate:"required"`
	Email        string             `json:"email" validate:"required,email"`
	Password     string             `json:"password" validate:"required"`
	AccountType  models.AccountType `json:"account_type" validate:"required"`
}

// LoginResponse represents the response data structure for the login API.
type LoginResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

// LoginResponseData represents the data section of the login response.
type LoginResponseData struct {
	JWT string `json:"jwt"`
}

type AuthenticateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateAccountInformation struct {
	Country        string `json:"country" validate:"required"`
	Manager        string `json:"manager" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required,numeric"`
	CompanyWebsite string `json:"company_website" validate:"required,url"`
}

type ResetPassword struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// ResetPassword resets the user's password using a JWT token.
//
// @Summary Reset password
// @Description Resets the user's password using a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ResetPassword true "New password"
// @Success 200 {string} string "Password reset successful"
// @Failure 400 {object} ErrorResponse
// @Router /auth/password-reset/new-password [post]
func (a *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*helpers.AuthTokenJwtClaim)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var input ResetPassword
	if err := c.BodyParser(&input); err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	// Update user's password
	user, userExist, err := a.userRepository.FindUserByCondition("email", claims.Email)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	if !userExist {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	if input.Password != input.ConfirmPassword {
		return helpers.Dispatch400Error(c, "password do not match", nil)
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	user.Password = hashedPassword
	_, err = a.userRepository.UpdateUserByCondition("email", claims.Email, user)
	if err != nil {
		return helpers.Dispatch400Error(c, "password update failed", err)
	}

	c.Status(http.StatusOK)
	return c.SendString("Password reset successful")
}

// VerifyOTPAndGenerateToken verifies the OTP and generates a JWT token for password reset.
//
// @Summary Verify OTP and generate JWT token
// @Description Verifies the provided OTP and generates a JWT token for password reset.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email path string true "User's email address"
// @Param otp path string true "One-time password (OTP)"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/verify-otp/{email}/{otp} [post]
func (a *AuthHandler) VerifyOTPAndGenerateToken(c *fiber.Ctx) error {
	email := c.Params("email")
	token := c.Params("otp")

	// Generate JWT token
	user, userExist, err := a.userRepository.FindUserByCondition("email", email)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	if !userExist {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	// Verify OTP
	valid := otp.OTPManage.VerifyOTP(user.Email, token)
	if !valid {
		return helpers.Dispatch400Error(c, "invalid OTP", nil)
	}

	jwtToken, err := helpers.GenerateToken(constant.JWTSecretKey, user.Email, user.FirstName+" "+user.LastName, user.UserId)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	c.Status(http.StatusOK)
	return c.JSON(LoginResponse{
		Success: true,
		Message: "OTP verified and token generated successfully",
		Data: LoginResponseData{
			JWT: jwtToken,
		},
	})
}

// SendOTPForPasswordReset sends an OTP to the provided email for password reset.
//
// @Summary Send OTP for password reset
// @Description Sends an OTP to the provided email address for password reset.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email query string true "User's email address"
// @Success 200 {string} string "OTP sent successfully"
// @Failure 400 {object} ErrorResponse
// @Router /auth/password-reset/send-otp [get]
func (a *AuthHandler) SendOTPForPasswordReset(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return fmt.Errorf("no email address provided")
	}

	// Check if the user exists
	user, userExist, err := a.userRepository.FindUserByCondition("email", email)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	if !userExist {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	// Generate and send OTP
	otpToken, err := otp.OTPManage.GenerateOTP(user.Email, time.Minute*10, 6)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	// Send OTP via email
	sendOTPEmail(user.LastName, user.Email, otpToken)

	c.Status(http.StatusOK)
	return c.SendString("OTP sent successfully")
}

func sendOTPEmail(name, email, token string) {
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}

	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10, 6)
	type OTP struct {
		Otp  string
		Name string
		Url  string
	}
	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	messageBody, err := helpers.ParseTemplateFile("account_reset.html", OTP{Otp: otpToken, Name: name})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	client := sendgrid.NewClient(constant.SendGridApiKey, constant.SenderEmail, "NturuCLI", "Reset Your Password", messageBody)
	err = client.Send(&to)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// Authenticate authenticates a user and generates a JWT token.
//
// @Summary Authenticate User
// @Description Authenticate a user by validating their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticateUser true "User credentials (email and password)"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/signin [post]
func (a *AuthHandler) Authenticate(c *fiber.Ctx) error {
	var input helpers.AuthenticateUser
	if err := c.BodyParser(&input); err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	user, userExist, err := a.userRepository.FindUserByCondition("email", input.Email)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	if !userExist {
		return helpers.Dispatch400Error(c, "invalid account credentials", nil)
	}

	if !user.EmailVerified {
		return helpers.Dispatch400Error(c, "account not verified", nil)
	}

	hashedPassword := []byte(user.Password)
	plainPassword := []byte(input.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)

	if err != nil {
		return helpers.Dispatch400Error(c, "email and password does not match", nil)
	}

	// update last login and ip
	time, err := helpers.TimeNow("Africa/Lagos")
	user.LastLogin = time
	user.IP = c.IP()

	if err != nil {
		return helpers.Dispatch400Error(c, err.Error(), err)
	}
	_, err = a.userRepository.UpdateUserByCondition("email", user.Email, user)
	if err != nil {
		return helpers.Dispatch400Error(c, err.Error(), err)
	}
	jwtToken, err := helpers.GenerateToken(constant.JWTSecretKey, user.Email, user.LastName+" "+user.FirstName, user.UserId)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "authenticated successfully",
		"data": map[string]string{
			"jwt": jwtToken,
		},
	})
}

func (a *AuthHandler) findUserOrError(email string) (user *models.User, err error) {
	user, userExist, err := a.userRepository.FindUserByCondition("email", email)
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// Register creates a new user account.
//
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body InputCreateUser true "User data to create an account"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/signup [post]
func (a *AuthHandler) Register(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	var input helpers.InputCreateUser
	if err := c.BodyParser(&input); err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	userFound, err := a.findUserOrError(input.Email)

	if userFound != nil && err == nil {
		return helpers.Dispatch500Error(c, fmt.Errorf("user already registered"))
	}

	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	// create record
	user := &models.User{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Password:      hash,
		UserId:        helpers.GenerateUUID(),
		IP:            c.IP(),
		Role:          models.UserRole,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		Status:        (models.InactiveAccount),
	}

	if err := a.userRepository.CreateUser(user); err != nil {
		return helpers.Dispatch500Error(c, err)
	}
	go sendVerificationEmail(user.FirstName, user.Email)

	c.Status(http.StatusCreated)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "register user successfully",
		"data": map[string]string{
			"id":    fmt.Sprint(user.UserId),
			"email": user.Email,
		},
	})

}

func sendVerificationEmail(name, email string) {
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}

	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10, 6)
	type OTP struct {
		Otp  string
		Name string
	}
	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	messageBody, err := helpers.ParseTemplateFile("verify_account.html", OTP{Otp: otpToken, Name: name})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
	client := sendgrid.NewClient(constant.SendGridApiKey, constant.SenderEmail, "NturuCLI", "Verify your email", messageBody)
	err = client.Send(&to)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// CompleteAccountInformation is a route handler that completes the account information for the authenticated user.
//
// This endpoint allows the authenticated user to update their account information.
//
// @Summary Complete account information
// @Description Completes the account information for the authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body UpdateAccountInformation true "Input data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /auth/complete [post]
// func (a *AuthHandler) CompleteAccountInformation(c *fiber.Ctx) error {
// 	claims, ok := c.Locals("claims").(*helpers.AuthTokenJwtClaim)
// 	if !ok {
// 		return fiber.ErrUnauthorized
// 	}
// 	var input helpers.UpdateAccountInformation

// 	if err := c.BodyParser(&input); err != nil {
// 		return helpers.Dispatch500Error(c, err)
// 	}
// 	user, _, err := a.userRepository.FindUserByCondition("email", claims.Email)

// 	if err != nil {
// 		return helpers.Dispatch400Error(c, "something went wrong", err.Error())
// 	}

// 	user.Country = input.Country
// 	user.PhoneNumber = input.PhoneNumber

// 	_, err = a.userRepository.UpdateUserByCondition("email", claims.Email, user)

// 	if err != nil {
// 		return helpers.Dispatch400Error(c, "something went wrong", err.Error())

// 	}
// 	return nil
// }
