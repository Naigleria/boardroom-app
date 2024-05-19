package main

import (
	//"fmt"
	"gorm/controllers"
	"gorm/initializers"
	//"time"

	
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	
	r := gin.Default()

	//lanzamon una funcion en un hilo paralelo utilizando goroutines para no efactar las ejecucion de los demas enpoints
	//ésta funcion verifica cada 60 segundos si hay alguna sala de juntas a la cual se le haya acabado el tiempo
	go controllers.StartPeriodicTask()

	//Operaciones CRUD para la sala de juntas
	r.POST("/api/boardroom/create_boardroom", controllers.CreateBoardroom)
	r.GET("/api/boardroom/", controllers.ReadAllBoardrooms)
	r.GET("/api/boardroom/:ID", controllers.ReadBoardroomById)
	r.PUT("/api/boardroom/:ID", controllers.UpdateBoardroomById)

	//Endpoint para reservar una sala de juntas
	r.POST("/api/boardroom/book_boardroom/:ID", controllers.BookBoardroomById)
	//Endpoint para liberar una sala de juntas
	r.GET("/api/boardroom/release_boardroom/:ID", controllers.ReleaseBoardroomById)

	r.Run()
}
