package controller

import (
	"app/config"
	"app/database"
	"app/modules/team/model"
	modell "app/modules/group/model"
	"app/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"os"
	"time"
)

type tempCreateTeamModel model.CreateTeamModel
type tempTeam model.Team
type updateGr model.UpdateTeamModel


// GetTeam Lấy danh sách tất cả team
// @Summary Get all Teams
// @Description Returns a list of all Teams
// @Tags Team
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var teams []model.Team
	results := database.DB.Select("*").Preload("Group").Order("team_id").Find(&teams)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = teams
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetAllTeam Lấy danh sách các team đã bị xoá
// @Summary Get all Teams (deleted)
// @Description Returns a list of all Teams (soft-deleted)
// @Tags Team
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team/all [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetAllTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var teams []model.Team
	results := database.DB.Select("*").Unscoped().Preload("Group").Where("deleted_at IS NOT NULL").Order("team_id").Find(&teams)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = teams
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetTeamByID returns information about a Team based on its ID
// @Summary Get a Team by ID
// @Description Returns information about a Team based on its ID
// @Tags Team
// @Accept json
// @Produce json
// @Param id path int true "ID of the Team"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team/{id} [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetTeamByID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var team model.Team
	results := database.DB.Select("*").Preload("Group").Where("team_id = ?", c.Params("id")).First(&team)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = team
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateTeam Tạo mới 1 team
// @Summary Create a new Team
// @Description Creates a new Team
// @Tags Team
// @Accept json
// @Produce json
// @Param body body []model.CreateTeamModel true "New Team information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team [post]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func CreateTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateTeamModel
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range payload {
		listCheck := []string{"TeamNameVN", "TeamNameEN", "TeamNameJP"}
		vItem := map[string]string{
			"TeamNameVN": item.TeamNameVN,
			"TeamNameEN": item.TeamNameEN,
			"TeamNameJP": item.TeamNameJP,
		}
		errors := utils.RequireCheck(listCheck, vItem, map[string]string{})

		if len(errors) > 0 {
			response.Status = false
			response.Message = config.GetMessageCode("MISSING_FIELDS")
			response.ValidateError = errors
			return c.JSON(response)
		}

		newTeam := model.Team{
			GroupID:  item.GroupID,
			TeamNameVN:   item.TeamNameVN,
			TeamNameEN:   item.TeamNameEN,
			TeamNameJP:   item.TeamNameJP,
			TeamShortcut: item.TeamShortcut,
			CreatedBy:     getUsername(c),
		}
		
		if err := database.DB.Create(&newTeam).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("CREATE_FAIL")
			return c.JSON(response)
		}
		
		tempNewTeam := tempTeam(newTeam)
		if err := tempNewTeam.AfterCreate(tx); err != nil {
			response.Status = false
			response.Message = config.GetMessageCode("AFTER_CREATE_FAIL")
			return c.JSON(response)
		}
		

	}

	response.Status = true
	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	return c.JSON(response)
}

// UpdateTeam cập nhật thông tin một Team dựa trên ID
// @Summary Cập nhật thông tin Team
// @Description Cập nhật thông tin một Team dựa trên ID
// @Tags Team
// @Accept json
// @Produce json
// @Param body body []model.UpdateTeamModel true "Thông tin Team cần cập nhật"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team [put]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func UpdateTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateTeamModel
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {
		listCheck := []string{"TeamNameVN", "TeamNameEN", "TeamNameJP"}
		vItem := map[string]string{"TeamNameVN": item.TeamNameVN, "TeamNameEN": item.TeamNameEN, "TeamNameJP": item.TeamNameJP}
		errors := utils.RequireCheck(listCheck, vItem, map[string]string{})

		if len(errors) > 0 {
			response.Message = "MISSING_FIELDS" //NOT_ID_EXISTS
			response.ValidateError = errors
			return c.JSON(response)
		}

		tempNTeam := updateGr(*item)
		
		if err := tempNTeam.BeforeUpdate(tx); err != nil {
			response.Status = false
			response.Message = config.GetMessageCode("ERROR_BEFORE_UPDATE")
			return c.JSON(response)
		}

		var team model.Team
		if item.TeamID != 0 {
			if item.IsDeleted {
				// Soft delete team
				results := db.First(&team, item.TeamID)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("NOT_ID_EXISTS")
					return c.JSON(response)
				}
				team.DeletedBy = getUsername(c)
				team.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			} else {
				// Update team
				results := db.First(&team, item.TeamID)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("NOT_ID_EXISTS")
					return c.JSON(response)
				}

				team.GroupID = item.GroupID
				team.TeamNameVN = item.TeamNameVN
				team.TeamNameEN = item.TeamNameEN
				team.TeamNameJP = item.TeamNameJP
				team.TeamShortcut = item.TeamShortcut


				if err := db.Model(&team).Updates(team).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}

			if err := db.Save(&team).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("SYSTEM_ERROR")
				return c.JSON(response)
			}
		} else {
			// Create team
			newTeam := model.Team{
				TeamNameVN:   item.TeamNameVN,
				TeamNameEN:   item.TeamNameEN,
				TeamNameJP:   item.TeamNameJP,
				TeamShortcut: item.TeamShortcut,
				CreatedBy:     getUsername(c),
			}

			if err := db.Create(&newTeam).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("CREATE_FAIL")
				return c.JSON(response)
			}
		}
	}

	tx.Commit()

	response.Status = true
	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	return c.JSON(response)
}

// DeleteTeam xóa một Team dựa trên ID
// @Summary Xóa Team
// @Description Xóa một Team dựa trên ID
// @Tags Team
// @Accept json
// @Produce json
// @Param id path int true "ID của Team"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team/{id} [delete]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func DeleteTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	teamID := c.Params("id")

	// Check if the team exists
	var team model.Team
	result := database.DB.First(&team, teamID)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}


	results := database.DB.First(&team, "team_id = ?", teamID)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Soft delete data: Update deleted_at field with the current time
	team.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	team.DeletedBy = getUsername(c)

	// Update data: update team by id from team table (database)
	if err := database.DB.Model(&team).Updates(&team).Error; err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}


	response.Status = true
	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	return c.JSON(response)
}

// RestoreTeam khôi phục một Team dựa trên ID
// @Summary Khôi phục Team
// @Description Khôi phục một Team dựa trên ID
// @Tags Team
// @Accept json
// @Produce json
// @Param id path int true "ID của Team"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /team/{id}/restore [put]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func RestoreTeam(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	teamID := c.Params("id")

	// Check if the team exists
	var team model.Team
	result := database.DB.Unscoped().First(&team, teamID)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}


	// Restore the team
	if err := database.DB.Unscoped().Model(&team).Update("deleted_at", nil).Error; err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("RESTORE_FAIL")
		return c.JSON(response)
	}

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)
}



// Hook được gọi sau khi thực hiện cập nhật (Update)
func (g *updateGr) BeforeUpdate(tx *gorm.DB) (err error) {
	if g.UpdatedBy == "1105" {
		fmt.Println(">>>>  it's updated by Admin....")
	}

	var group modell.Group
	results := database.DB.Select("*").Where("group_id = ?", g.GroupID).First(&group)
	if results.Error != nil {
		return results.Error
	}

	if !tx.Statement.Changed("team_id", "team_name_vn", "team_name_en", "team_name_jp", "team_shortcut") {
		WriteLog("DATA_NOT_CHANGED")
	}

	if tx.Statement.Changed("LogVersion") {
		return errors.New(config.Config("DATA_NOT_CHANGED"))
	}

	team := tx.Statement.Dest.(*model.Team)
	team.LogVersion++


	return nil
}



func (g *tempTeam) AfterCreate(tx *gorm.DB) (err error) {
    if g.CreatedBy == "1105" {
        fmt.Println(">>>>  it's created by Admin....")
		WriteLog(">>>>  it's created by Admin....")
    }
    return nil
}


func WriteLog(logEntry string) error {
	// Tạo một timestamp để thêm vào log
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntryWithTimestamp := fmt.Sprintf("[%s] %s\n", timestamp, logEntry)

	// Mở file log.txt trong chế độ append
	file, err := os.OpenFile("assets/log/log.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ghi log vào file
	_, err = file.WriteString(logEntryWithTimestamp)
	if err != nil {
		return err
	}

	return nil
}

func getUsername(c *fiber.Ctx) string {
	response := new(config.DataResponse)

	tokenData, err := utils.ExtractTokenData(c)
	if err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("ERROR_GET_USERNAME")
	}

	return tokenData.Username
}
