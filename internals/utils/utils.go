package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	ID      string `json:"id"`
	Company string `json:"company"`
	Branch  string `json:"branch"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token with the given user ID, company ID, branch ID, and role.
// It sets the expiration time to 24 hours from the current time.
// It returns the generated token as a string and any error encountered.
func GenerateJWT(id, company, branch, role string) (string, error) {
	// Set the expiration time to 24 hours from the current time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create a Claims object with the given user ID, company ID, branch ID, role, and expiration time
	claims := &Claims{
		ID:      id,
		Company: company,
		Branch:  branch,
		Role:    role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new JWT token with the claims object and the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the JWT key and return it as a string
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token. It returns the claims in the token and any error encountered.
func ParseJWT(tokenString string) (*Claims, error) {
	// Create a new instance of the Claims struct that will hold the token claims
	claims := &Claims{}

	// Parse the token string into a JWT token, validating its signature and setting the claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Return the JWT key used to sign the token
		return jwtKey, nil
	})

	// If there was an error parsing or validating the token, return it
	if err != nil {
		return nil, err
	}

	// If the token is invalid, return an error
	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	// Return the claims in the token
	return claims, nil
}
