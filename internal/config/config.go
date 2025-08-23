package config

import "fmt"

type DatabaseConfig struct{
	Host string
	Port int
	User string
	Password string
	DatabaseName string
	SSLMode string
}

func LoadDBConfig() *DatabaseConfig{
	return &DatabaseConfig{
		Host:"localhost",
		Port:5433,
		User:"postgres",
		Password: "password",
		DatabaseName: "userdb",
		SSLMode: "disable",
	}
}
func (cfg *DatabaseConfig) GetConnectionString()string{
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
cfg.Host,cfg.Port,cfg.User,cfg.Password,cfg.DatabaseName,cfg.SSLMode)
}
