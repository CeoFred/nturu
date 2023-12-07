package helpers

import (
	// "github.com/nturu/microservice-template/internal/models"
	"github.com/golang-jwt/jwt"
)

type InputCreateUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type UpdateAccountInformation struct {
	Country        string `json:"country" validate:"required"`
	Manager        string `json:"manager" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	CompanyWebsite string `json:"company_website"`
}
type ListProjectOffsets struct {
	Projects []string `json:"projects" validate:"required"`
}

type OtpVerify struct {
	Token string `json:"token" validate:"required,len=5"`
	Email string `json:"email" validate:"required,email"`
}

type AccountReset struct {
	Email string `json:"email" validate:"required,email"`
}
type AuthenticateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type IError struct {
	Field string
	Tag   string
	Value string
}

type AuthTokenJwtClaim struct {
	Email  string
	Name   string
	UserId string
	jwt.StandardClaims
}
type ProjectType string

type InputCreateProject struct {
	Country              string      `json:"country" validate:"required"`
	LaunchDate           string      `json:"launch_date" validate:"required"`
	ActiveDaysPerWeek    int         `json:"active_days_per_week" validate:"required"`
	Category             string      `json:"category" validate:"required"`
	Type                 ProjectType `json:"type" validate:"required"`
	LifeSpan             int         `json:"life_span"`
	Capacity             int         `json:"capacity"`
	CapacityUtility      int         `json:"capacity_utility"`
	Size                 int         `json:"size"`
	NewTrees             int         `json:"new_trees"`
	LifeSpanPerRefill    int         `json:"life_span_per_refill"`
	LifeSpanPerCookStove int         `json:"life_span_per_cook_stove"`
	WasteAmount          int         `json:"waste_amount"`
	SiteAddress          string      `json:"site_address" validate:"required"`
	ActiveHoursPerDay    int         `json:"active_hours_per_day" validate:"`
	Description          string      `json:"description" validate:"required"`
	ImageUrl             string      `json:"image_url" validate:"required"`
}

type AccountStatus int
