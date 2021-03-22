package structs

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func init() {
	fmt.Println("package: structs.jwt - initialized")
}

// Keys struct
type Keys struct {
	Keys []Key `json:"keys"`
}

// Key struct
type Key struct {
	Kid string `json:"kid"`
	Nbf int    `json:"nbf"`
	Use string `json:"use"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
}

// UserClaims struct
type UserClaims struct {
	jwt.StandardClaims
	Username string
}
