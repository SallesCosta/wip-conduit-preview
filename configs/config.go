package configs

import (
	"log"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"

	"os"
)

type Conf struct {
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	WebServePort string `mapstructure:"WEB_SERVER_PORT"`
	JwtExpiresIn int    `mapstructure:"JWT_EXPIRESIN"`
	TokenAuth    *jwtauth.JWTAuth

	DBDriver string `mapstructure:"DB_DRIVER"`
}

func LoadConfig() (*Conf, error) {
	if err := godotenv.Load("cmd/server/.env"); err != nil {
		panic(err)
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "8080"
	}

	jwtExpiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRESIN"))
	if err != nil {
		log.Fatalf("Erro ao converter JWT_EXPIRESIN para inteiro: %v", err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	return &Conf{
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       port,
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		WebServePort: os.Getenv("WEB_SERVER_PORT"),
		JwtExpiresIn: jwtExpiresIn,
		TokenAuth:    tokenAuth,
	}, err
}
