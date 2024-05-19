package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

//metodo para cargar las variables de entorno
func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
