package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/nturu/microservice-template/constants"

	"github.com/nturu/microservice-template/internal/helpers"
	"github.com/nturu/microservice-template/internal/models"
	"github.com/nturu/microservice-template/internal/repository"
)

var (
	env_ = constants.New()
)

type UserHandler struct {
	userRepository *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepo,
	}
}

type FileUploadResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserProfile struct {
	Email       string                   `json:"email"`
	Userid      string                   `json:"userid"`
	Role        models.AccountPermission `json:"role"`
	Country     string                   `json:"country"`
	PhoneNumber string                   `json:"phone_number"`
	Status      string                   `json:"status"`
	CreatedAt   time.Time                `json:"created_at"`
	AccountType models.AccountType       `json:"account_type"`
}

type UpdateUserProfileInput struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// UpdateUserProfile is a route handler that handles updating the user profile
//
// # This endpoint is used to update the user profile
//
// @Summary Update user profile
// @Description Updates some details about the user
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body UpdateUserProfileInput true "update user profile"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /user [put]
func (u *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	var input UpdateUserProfileInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.Dispatch500Error(c, err)
	}

	claims, ok := c.Locals("claims").(*helpers.AuthTokenJwtClaim)
	if !ok {
		return fiber.ErrUnauthorized
	}
	user, found, err := u.userRepository.FindUserByCondition("email", claims.Email)

	if err != nil {
		return helpers.Dispatch400Error(c, "something went wrong", err)
	}
	if !found {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	user.PhoneNumber = input.PhoneNumber

	_, err = u.userRepository.UpdateUserByCondition("email", user.Email, user)
	if err != nil {
		return helpers.Dispatch400Error(c, "something went wrong", err)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"success": true,
		"message": "updated successfully",
	})
}

// UserProfile is a route handler that retrieves the user profile of the authenticated user.
//
// This endpoint is used to get the profile information of the authenticated user based on the JWT claims.
//
// @Summary Get user profile
// @Description Retrieves the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfile
// @Failure 401 {object} ErrorResponse
// @Router /user/profile [get]
func (u *UserHandler) UserProfile(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*helpers.AuthTokenJwtClaim)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, found, err := u.userRepository.FindUserByCondition("email", claims.Email)

	if err != nil {
		return helpers.Dispatch400Error(c, "something went wrong", err)
	}
	if !found {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	p := UserProfile{
		Email:       user.Email,
		Userid:      user.UserId,
		Role:        user.Role,
		Country:     user.Country,
		PhoneNumber: user.PhoneNumber,
		Status:      string(models.AccountStatus(user.Status)),
		CreatedAt:   user.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    p,
	})
}

// FileUpload is a route handler that handles file uploads.
//
// This endpoint is used to upload a file.
//
// @Summary Upload a file
// @Description Handles file uploads
// @Tags User
// @Accept mpfd
// @Produce json
// @Param image formData file true "File to upload"
// @Success 200 {object} FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/file-upload [post]
func (u *UserHandler) FileUpload(c *fiber.Ctx) error {

	file, err := c.FormFile("image")
	if err != nil {

		return helpers.Dispatch400Error(c, err.Error(), err)
	}
	resp, err := uploadFile(file)
	if err != nil {
		return helpers.Dispatch400Error(c, err.Error(), err)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"success": true,
		"message": "uploaded successfully",
		"data": map[string]interface{}{
			"url": resp.SecureURL,
		},
	})

}

func uploadFile(file *multipart.FileHeader) (resp *uploader.UploadResult, err error) {

	// Open the uploaded file
	fileOpened, err := file.Open()
	if err != nil {
		// Handle error
		return nil, err
	}

	defer fileOpened.Close()
	url := fmt.Sprintf("cloudinary://%s:%s@%s", env_.CloudinaryAPIKey, env_.CloudinaryApiSecret, env_.CloudinaryName)

	cld, err := cloudinary.NewFromURL(url)

	if err != nil {
		return nil, err
	}
	var ctx = context.Background()
	// Upload the image to Cloudinary
	resp, err = cld.Upload.Upload(ctx, fileOpened, uploader.UploadParams{PublicID: file.Filename,
		Folder: "NturuCLI",
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// @Summary Company Logo Upload
// @Description Company logo upload
// @Tags User
// @Accept json
// @Produce json
// @Param requestBody body interface{} true "Upload a company logo"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /user/logo [post]
func (u *UserHandler) UploadLogo(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*helpers.AuthTokenJwtClaim)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, found, err := u.userRepository.FindUserByCondition("email", claims.Email)

	if err != nil {
		return helpers.Dispatch400Error(c, "something went wrong", err)
	}
	if !found {
		return helpers.Dispatch400Error(c, "user not found", nil)
	}

	file, err := c.FormFile("logo")
	if err != nil {
		// Handle error
		return helpers.Dispatch400Error(c, err.Error(), err)
	}

	// Open the uploaded file
	resp, err := uploadFile(file)
	if err != nil {

		return helpers.Dispatch400Error(c, err.Error(), err)
	}

	_, err = u.userRepository.UpdateUserByCondition("email", user.Email, user)
	if err != nil {
		return helpers.Dispatch400Error(c, "something went wrong", err)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"success": true,
		"message": "uploaded successfully",
		"data": map[string]interface{}{
			"url": resp.SecureURL,
		},
	})

}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}
