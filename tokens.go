package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

func SetJWTKey(key string) {
	if key == "" {
		panic("JWT_SECRET_KEY cannot be empty")
	}
	jwtKey = []byte(key)
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}


func GenerateJWT(email string, duration time.Duration) (string, error) {
	now := time.Now()
	expirationTime := now.Add(duration)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	if IsTokenRevoked(tokenStr) {
		return nil, errors.New("token is revoked")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func RevokeToken(tokenStr string) {
	if IsTokenRevoked(tokenStr) {
		return
	}

	tempClaims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, tempClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	var tokenExpiryTime time.Time
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 && tempClaims.ExpiresAt != nil {
				tokenExpiryTime = tempClaims.ExpiresAt.Time
			} else {
				return
			}
		} else if tempClaims.ExpiresAt != nil {
			tokenExpiryTime = tempClaims.ExpiresAt.Time
		} else {
			return
		}
	} else if tempClaims.ExpiresAt != nil {
		tokenExpiryTime = tempClaims.ExpiresAt.Time
	} else {
		return
	}

	revokedEntry := RevokedToken{
		TokenString: tokenStr,
		ExpiresAt:   tokenExpiryTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = RevokedTokenCollection.InsertOne(ctx, revokedEntry)
}

func IsTokenRevoked(tokenStr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var foundToken RevokedToken
	err := RevokedTokenCollection.FindOne(ctx, bson.M{"tokenString": tokenStr}).Decode(&foundToken)
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil {	
		return false
	}
	return true
}
