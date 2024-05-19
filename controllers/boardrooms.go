package controllers

import (
	"fmt"
	"gorm/initializers"
	"gorm/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//funcion que crea una nueva sala de juntas
func CreateBoardroom(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Capacity uint   `json:"capacity"`
	}

	if c.ShouldBindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	boardroom := models.Boardroom{
		Name: body.Name,
		Capacity:    body.Capacity,
	}

	if err := initializers.DB.Create(&boardroom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create boardroom",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Boardroom created sucessfully",
	})
}

//funcion para leer todas las salas de juntas existentes
func ReadAllBoardrooms(c *gin.Context){
	boardrooms := models.Boardrooms{}
	initializers.DB.Find(&boardrooms)
	c.JSON(http.StatusOK, boardrooms)
}

//funcion para leer solo una sala de juntas mediante su ID
func ReadBoardroomById(c *gin.Context){
	boardroomId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
		return
	}

	var boardroom models.Boardroom
	err = initializers.DB.First(&boardroom, boardroomId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boardroom not found"})
		return
	}

	c.JSON(http.StatusOK, boardroom)
}

//funcion para actualizar una sala de juntas mediante ID
func UpdateBoardroomById(c *gin.Context){
	boardroomId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
		return
	}

	var boardroom models.Boardroom
	err = initializers.DB.First(&boardroom, boardroomId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boardroom not found"})
		return
	}

	var boardroom_new models.Boardroom
	if err := c.ShouldBindJSON(&boardroom_new); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON"})
		return
	}

	boardroom_new.ID = uint(boardroom.ID)

	if err := initializers.DB.Save(&boardroom_new).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error trying to update boardroom"})
		return
	}

	// Response with updated boardroom
	c.JSON(http.StatusOK, boardroom_new)
}

//funcion para eliminar una sala de juntas mediante ID
func DeleteBoardroomById(c *gin.Context){
	boardroomId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
		return
	}

	var boardroom models.Boardroom
	err = initializers.DB.First(&boardroom, boardroomId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boardroom not found"})
		return
	}

	err = initializers.DB.Delete(&boardroom).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error trying to delete boardroom"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Boardroom deleted successfully"})
}

//funcion para reservar una sala de juntas recibiendo un horario inicial y final
func BookBoardroomById(c *gin.Context){
	boardroomId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
		return
	}

	type Boardroom struct {
		ID uint `gorm:"column:boardroom_id;primaryKey" json:"boardroom_id"`
		Available bool `json:"available"`
	}

	var boardroom Boardroom
	result:= initializers.DB.Debug().Select("boardroom_id, available").Where("boardroom_id = ?", boardroomId).First(&boardroom)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to get boardroom availability",
		})
		return
	}
	
	fmt.Println(boardroom)
	//validamos aqui si la sala de juntas está disponible, si no mandamos un mensaje de error
	if boardroom.Available == false{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unavailable boardroom",
		})
		return
	}

	
	var body struct {
		InitialSchedule string `json:"initial_schedule"`
		FinalSchedule   string `json:"final_schedule"`
	}

	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	const layout ="2006-01-02 15:04"

	parsedInitial, err:= time.Parse(layout, body.InitialSchedule)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid format for initial schedule"})
		return
	}
	parsedFinal, err:= time.Parse(layout, body.FinalSchedule)

	

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid format for final schedule"})
		return
	}

	// Calculamos la duración de la reserva
	duration := parsedFinal.Sub(parsedInitial)

	//aqui validamos que la reserva de la sala no exceda las 2 horas (120 minutos)
	//si excede las 2 horas mandamos un mensaje de error
	if duration.Minutes() > 120 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation cannot exceed 2 hours"})
		return
	}

	values:= models.Boardroom{ID: uint(boardroomId), InitialSchedule: body.InitialSchedule, FinalSchedule: body.FinalSchedule, Available: false}

	
	if err = initializers.DB.Model(&values).Select("initial_schedule", "final_schedule", "available").Updates(values).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to book boardroom",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Boardroom booked successfully",
	})
}

//funcion para liberar una sala de juntas manualmente mediante ID
func ReleaseBoardroomById(c *gin.Context){
	boardroomId, err := strconv.Atoi(c.Param("ID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
		return
	}

	values:= models.Boardroom{ID: uint(boardroomId),  InitialSchedule:"", FinalSchedule:"", Available: true}

	if err = initializers.DB.Debug().Model(&values).Updates(map[string]interface{}{"initial_schedule": values.InitialSchedule, "final_schedule": values.FinalSchedule, "Available": values.Available}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to release boardroom",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Boardroom released successfully",
	})
}

//funcion para determinar si alguna sala de juntas ya se le acabó el tiempo, si es asi la ponemos como disponible	
func StartPeriodicTask() {
    ticker := time.NewTicker(60 * time.Second) // el ticker nos permite determinar cada cuanto tiempo se ejecutará la verificacion
    defer ticker.Stop()

	//iniciamos un bucle infinito para estar checando constantemente si ya se terminó el tiempo de alguna sala
    for {
        select {
        case <-ticker.C:
			
			boardrooms := models.Boardrooms{}
			initializers.DB.Find(&boardrooms)
			
			//nos traemos todas las salas y verificamos el tiempo final de cada una de ellas
			for _, element := range boardrooms{
				
				
				if element.FinalSchedule != ""{
					now:=time.Now()
					const layout ="2006-01-02 15:04"

					formattedTimeNow := now.Format(layout)

					formattedTime, err := time.Parse(layout, formattedTimeNow)
					if err != nil {
						fmt.Println("Error formatting current time:", err)
						return
					}
					
					var final time.Time
					final, err = time.Parse(layout, element.FinalSchedule)

					if err != nil {
						fmt.Println("Error formatting final time:", err)
						return
					}

					//al tiempo actual le restamos el horario de reserva final de la sala
					duration:=formattedTime.Sub(final)
					//si sobra tiempo quiere decir que ya se le terminó el tiempo y ponemos la sala en disponible
					if duration.Minutes() > 0 {
						//liberamos sala de juntas
						values:= models.Boardroom{ID: element.ID,  InitialSchedule:"", FinalSchedule:"", Available: true}
	
						if err = initializers.DB.Debug().Model(&values).Updates(map[string]interface{}{"initial_schedule": values.InitialSchedule, "final_schedule": values.FinalSchedule, "Available": values.Available}).Error; err != nil {
							fmt.Println("Failed to release boardroom timer")
							return
						}
	
						fmt.Println("boardroom released automatically!")
					}
				}
				
			}
			fmt.Println("Executed function: ")
        }
    }
}