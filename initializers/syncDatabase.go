package initializers

import "gorm/models"

//método para migrar los modelos de la carpeta "models" a la base de datos
func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Boardroom{})
}
