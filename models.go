package main

import (
	"context"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	HashedPassword []byte             `bson:"hashedPassword"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}


type RevokedToken struct {
	TokenString string    `bson:"tokenString"`
	ExpiresAt   time.Time `bson:"expiresAt"`
}

func AddUser(email string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingUser := UserCollection.FindOne(ctx, bson.M{"email": email})
	if existingUser.Err() == nil {
		return ErrUserExists
	}
	if existingUser.Err() != mongo.ErrNoDocuments {
		return existingUser.Err()
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := User{
		Email:          email,
		HashedPassword: hashedPwd,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err = UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func AuthenticateUser(email string, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	return err == nil
}

func CheckUserExists(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return err == nil
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
