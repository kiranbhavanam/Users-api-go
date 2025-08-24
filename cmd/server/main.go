package main

import (
	"fmt"
	"net/http"
	"user-management/internal/config"
	"user-management/internal/handlers"
	"user-management/internal/middleware"
	"user-management/internal/repository"
	"user-management/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main(){
	cfg:=config.LoadConfig()
	connectionstring:=config.LoadDBConfig().GetConnectionString()
	repo,err:=repository.NewPostgresRepository(connectionstring)
	if err!=nil{
		fmt.Println("error while loading db",err)
	}
	service:=service.NewUserService(cfg,repo)
	handler:=handlers.NewUserHandler(service)

	router:=mux.NewRouter()
	router.HandleFunc("/users",handler.CreateHandler).Methods("POST")
	router.HandleFunc("/auth/login",handler.LoginHandler).Methods("POST")
	//protected routes authenticationrequired
	protected:=router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	
	protected.HandleFunc("/users",handler.GetAllHandler).Methods("GET")
	protected.HandleFunc("/users/{id:[0-9]+}",handler.GetByIDHandler).Methods("GET")
	protected.HandleFunc("/users/{id:[0-9]+}",handler.UpdateHandler).Methods("PUT")
	protected.HandleFunc("/users/{id:[0-9]+}",handler.DeleteHandler).Methods("DELETE")

	fmt.Println("starting user server on 8080:")
	if err:=http.ListenAndServe(":8080",router);err!=nil{
		fmt.Println("unable to start server")
	}

}