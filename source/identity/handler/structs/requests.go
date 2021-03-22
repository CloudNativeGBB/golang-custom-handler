package structs

import (
	"fmt"
)

func init() {
	fmt.Println("package: structs.requests - initialized")
}

// CreateUserRequest struct
type CreateUserRequest struct {
	GivenName string `json:"givenName"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// AuthUserRequest struct
type AuthUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
