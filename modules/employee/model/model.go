package model

import (
	"time"
	"gorm.io/gorm"
	"app/modules/group/model"
)

type Model struct {
	ID        uint `gorm:"primarykey;column:team_id;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


type Team struct {
	Model
	GroupID  int               `gorm:"column:group_id;not null"`
	Group    model.Group  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TeamNameVN   string            `gorm:"column:team_name_vn;size:100;not null"`
	TeamNameEN   string            `gorm:"column:team_name_en;size:100;not null"`
	TeamNameJP   string            `gorm:"column:team_name_jp;size:100;not null"`
	TeamShortcut string            `gorm:"column:team_shortcut;size:5"`
	LogVersion    int64             `gorm:"column:log_version;default:0"`
	CreatedBy  string     `gorm:"column:created_by;size:15"`
	UpdatedBy  string     `gorm:"column:updated_by;size:15"`	
	DeletedBy  string     `gorm:"column:deleted_by;size:15"`
}


type CreateTeamModel struct {
	GroupID  int    `json:"group_id" validate:"required"`
	TeamNameVN   string `json:"team_name_vn" validate:"required"`
	TeamNameEN   string `json:"team_name_en" validate:"required"`
	TeamNameJP   string `json:"team_name_jp" validate:"required"`
	TeamShortcut string `json:"team_shortcut"`
	CreatedBy     string `json:"created_by"`
}

type UpdateTeamModel struct {
	TeamID       int    `json:"team_id" validate:"required"`
	GroupID  int    `json:"group_id" validate:"required"`
	TeamNameVN   string `json:"team_name_vn" validate:"required"`
	TeamNameEN   string `json:"team_name_en" validate:"required"`
	TeamNameJP   string `json:"team_name_jp" validate:"required"`
	TeamShortcut string `json:"team_shortcut"`
	IsDeleted     bool   `json:"is_deleted"`
	UpdatedBy     string `json:"updated_by"`
}

// type HeaderTeam struct {
// 	CreatedBy  string     `gorm:"column:created_by;size:15"`
// 	UpdatedBy  string     `gorm:"column:updated_by;size:15"`	
// 	DeletedBy  string     `gorm:"column:deleted_by;size:15"`
// }



// Tên bảng trong CSDL
func (Team) TableName() string {
	return "tbl_team"
}
