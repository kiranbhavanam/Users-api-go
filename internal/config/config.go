package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct{
	dburl string
	JWTSecret string
	JWTExpiry time.Duration
}

type DatabaseConfig struct{
	Host string
	Port int
	User string
	Password string
	DatabaseName string
	SSLMode string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifeTime time.Duration
}
func LoadConfig() *Config{
	dbURL:=LoadDBConfig().GetConnectionString()
	if dbURL==""{
		log.Fatal("DB url can't be empty ")
	}
	secretkey:=getEnv("JWTSecret","secret_key")
	expiry,_:=time.ParseDuration(getEnv("JWTExpiry","5m"))
	return &Config{
		dburl: dbURL,
		JWTSecret: secretkey,
		JWTExpiry: expiry,
	}
}
func LoadDBConfig() *DatabaseConfig{
	port,_:=strconv.Atoi(getEnv("DB_PORT","5433"))
	return &DatabaseConfig{
		Host:getEnv("DB_HOST","localhost"),
		Port:port,
		User:getEnv("DB_USER","postgres"),
		Password: getEnv("DB_PASSWORD","password"),
		DatabaseName: getEnv("DB_NAME","userdb"),
		SSLMode: getEnv("DB_SSLMODE","disable"),
		MaxOpenConns:25,
		MaxIdleConns:25,
		MaxLifeTime:5*time.Minute,
	}
}
func (cfg *DatabaseConfig) GetConnectionString()string{
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
cfg.Host,cfg.Port,cfg.User,cfg.Password,cfg.DatabaseName,cfg.SSLMode)
}

func getEnv(value string,def string)string{
	if val:=os.Getenv(value);val!=""{
		return val
	}
	return def
}