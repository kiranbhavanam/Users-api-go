package auth

import (
	"errors"
	"fmt"
	"time"
	"user-management/internal/config"
	"user-management/internal/model"

	"github.com/golang-jwt/jwt/v5"
)
func GenerateToken(u *model.User,cfg *config.Config)(string,error){
	claims:=&model.AccessClaims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprintf("%d",u.ID),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiry)),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS512,claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string,cfg *config.Config)(*model.AccessClaims,error){
	claims:=&model.AccessClaims{}

	token,err:=jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(interface{},error){
		return []byte(cfg.JWTSecret),nil
	})
	if err!=nil{
		return nil,err
	}
	if !token.Valid{
		return nil,errors.New("invalid token")
	}
	return claims,nil
}