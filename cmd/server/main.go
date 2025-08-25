package main

import (
	"log/slog"
	"net/http"
	"os"
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
	loghandler:=slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: false,
	})
	slog.SetDefault(slog.New(loghandler))
	slog.Info("application starting", "version", "1.0.0")

	connectionstring:=config.LoadDBConfig().GetConnectionString()
	repo,err:=repository.NewPostgresRepository(connectionstring)
	if err!=nil{
		slog.Error("error while loading db","error",err)
		os.Exit(1)
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

	slog.Info("starting user server","port", 8080)
	if err:=http.ListenAndServe(":8080",router);err!=nil{
		slog.Error("unable to start server","error",err,"port",8080)
	}

}