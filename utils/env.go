package utils

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
)

type Vars struct {
	SECRET_KEY  string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	DB_SSLMODE  string
	DB_TIMEZONE string
}

type Environment struct {
	envFile string
}

func NewEnvironment(envFile string) *Environment {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %v file", envFile)
		panic(err)
	}
	return &Environment{
		envFile: envFile,
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func (env *Environment) GetEnv(key string) string {
	return os.Getenv(key)
}

func (env *Environment) GetEnvs() *Vars {
	return &Vars{
		SECRET_KEY:  os.Getenv("SECRET_KEY"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE: os.Getenv("DB_TIMEZONE"),
	}
}

func (env *Environment) Format(s string, v interface{}) string {
	t, b := new(template.Template), new(strings.Builder)
	template.Must(t.Parse(s)).Execute(b, v)
	return b.String()
}
