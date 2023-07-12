package model

import (
	"time"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primarykey;column:shift_id;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


type Shift struct {
	Model
	EmployeeID uint `gorm:"
}

type ShiftChild struct {
	timeStart time.Time  `gorm:"column:time_start;not null"`
	timeEnd time.Time `gorm:"column:time_end;not null"`
	breakStart time.Time `gorm:"column:break_start;not null"`
	breakEnd time.Time `gorm:"column:break_end;not null"`
	OverTime time.Time `gorm:"column:over_time;not null"`
	checkStatus bool `gorm:"column:check_status;default:false"`
}



// Tên bảng trong CSDL
func (Shift) TableName() string {
	return "tbl_shift"
}


