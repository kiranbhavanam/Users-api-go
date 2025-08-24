package model

import (
		"github.com/golang-jwt/jwt/v5"

)
type User struct{
	ID int 	`json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name,omitempty"`
	IsActive bool `json:"isactive,omitempty"`
}

type AccessClaims struct{
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct{
	User User  `json:"user"`
	Token string `json:"token"`
	Message string `json:"message"`
}