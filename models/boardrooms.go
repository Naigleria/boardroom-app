package models


import ()


//modelo para una sala de juntas
type Boardroom struct {
    ID uint   `gorm:"column:boardroom_id;primaryKey;autoIncrement" json:"boardroom_id"`
    InitialSchedule    string   `json:"initial_schedule"`
    FinalSchedule      string   `json:"final_schedule"`
    Available          bool    	`gorm:"not null;default:true"`
    Name               string   `json:"name"`
    Capacity           uint     `json:"capacity"`
    
}

type Boardrooms []Boardroom

