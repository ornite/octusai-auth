package models

import (
	"regexp"

	"github.com/go-playground/validator"
)

type User struct {
	ID string `bson:"_id"`    
	Username string `bson:"username" validate:"min=3"`           
	Email    string `bson:"email" validate:"email"`           
	Password string `bson:"password" validate:"password"`       
}

// Validator instance
var validate *validator.Validate

func init() {
	// Initialize the validator
	validate = validator.New()
	// Register custom validation: password
	validate.RegisterValidation("password", validatePassword)
}

// validatePassword verifies if the password meets the specified criteria: at least one uppercase, one lowercase, one digit, and min 9 characters.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Regex to check the password complexity
	var (
		uppercase = regexp.MustCompile(`[A-Z]`)
		lowercase = regexp.MustCompile(`[a-z]`)
		number    = regexp.MustCompile(`[0-9]`)
	)
	if len(password) < 9 || !uppercase.MatchString(password) || !lowercase.MatchString(password) || !number.MatchString(password) {
		return false
	}
	return true
}

func (u *User) Validate() error {
	// Use the validator to validate the struct
	return validate.Struct(u)
}
