package model

type User struct{
	ID int 	`json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name,omitempty"`
	IsActive bool `json:"isactive,omitempty"`
}