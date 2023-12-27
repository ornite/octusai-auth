package services

import (
	"context"
	"fmt"
	"time"

	models "auth/src/models"
	utils "auth/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService represents the service for managing user-related operations.
type UserService struct {
	collection *mongo.Collection // MongoDB collection for users
}

// NewUserService creates a new UserService.
func NewUserService(db *mongo.Database, collectionName string) *UserService {
	collection := db.Collection(collectionName)
	return &UserService{collection}
}

// RegisterUser registers a new user.
func (us *UserService) RegisterUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate the user fields.
    if err := user.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err) // Return the validation error directly
    }

	// Check if a user with the same email already exists
	existingUser := &models.User{}
	filter := bson.M{"email": user.Email}

	err := us.collection.FindOne(ctx, filter).Decode(existingUser)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update the user's password with the hashed version
	user.Password = string(hashedPassword)

	// Insert the user into the database
	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

// LoginUser performs user authentication and returns a user if successful.
func (us *UserService) LoginUser(email, password string) (*models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the user by email
	filter := bson.M{"email": email}

	var user models.User
	err := us.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", fmt.Errorf("user not found")
		}
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}

	// Compare the stored hashed password with the input password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("password does not match")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	// Return the user and the token
	return &user, token, nil
}
