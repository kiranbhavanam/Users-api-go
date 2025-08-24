package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-management/internal/auth"
	"user-management/internal/config"
)

func JWTMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		authHeader:=r.Header.Get("Authorization")
		if authHeader==""{
			http.Error(w,`{"error":"Authorization header required"}`,http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader,"Bearer "){
			http.Error(w,`{"error":"Invalid auth format. Use: Bearer <token>}`,http.StatusUnauthorized)
			return
		}

		tokenString:=strings.TrimPrefix(authHeader,"Bearer ")
		claims,err:=auth.ValidateToken(tokenString,config.LoadConfig())
		if err!=nil{
			http.Error(w,`{"error":"Invalid or expired token}`,http.StatusUnauthorized)
			return
		}
		ctx:=context.WithValue(r.Context(),"claims",claims)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}