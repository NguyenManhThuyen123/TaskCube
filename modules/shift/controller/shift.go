package controller

import (
	"app/config"
	"app/database"
	"app/modules/shift/model"
	modell "app/modules/department/model"
	"app/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"os"
	"time"
)

type tempCreateGroupModel model.CreateGroupModel
type tempGroup model.Group
type updateGr model.UpdateGroupModel


// GetGroup Lấy danh sách tất cả group
// @Summary Get all Groups
// @Description Returns a list of all Groups
// @Tags Group
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var groups []model.Group
	results := database.DB.Select("*").Preload("Department").Order("group_id").Find(&groups)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = groups
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetAllGroup Lấy danh sách các group đã bị xoá
// @Summary Get all Groups (deleted)
// @Description Returns a list of all Groups (soft-deleted)
// @Tags Group
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group/all [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetAllGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var groups []model.Group
	results := database.DB.Select("*").Unscoped().Preload("Department").Where("deleted_at IS NOT NULL").Order("group_id").Find(&groups)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = groups
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetGroupByID returns information about a Group based on its ID
// @Summary Get a Group by ID
// @Description Returns information about a Group based on its ID
// @Tags Group
// @Accept json
// @Produce json
// @Param id path int true "ID of the Group"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group/{id} [get]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func GetGroupByID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var group model.Group
	results := database.DB.Select("*").Preload("Department").Where("group_id = ?", c.Params("id")).First(&group)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = group
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateGroup Tạo mới 1 group
// @Summary Create a new Group
// @Description Creates a new Group
// @Tags Group
// @Accept json
// @Produce json
// @Param body body []model.CreateGroupModel true "New Group information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group [post]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func CreateGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateGroupModel
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range payload {
		listCheck := []string{"GroupNameVN", "GroupNameEN", "GroupNameJP"}
		vItem := map[string]string{
			"GroupNameVN": item.GroupNameVN,
			"GroupNameEN": item.GroupNameEN,
			"GroupNameJP": item.GroupNameJP,
		}
		errors := utils.RequireCheck(listCheck, vItem, map[string]string{})

		if len(errors) > 0 {
			response.Status = false
			response.Message = config.GetMessageCode("MISSING_FIELDS")
			response.ValidateError = errors
			return c.JSON(response)
		}

		newGroup := model.Group{
			DepartmentID:  item.DepartmentID,
			GroupNameVN:   item.GroupNameVN,
			GroupNameEN:   item.GroupNameEN,
			GroupNameJP:   item.GroupNameJP,
			GroupShortcut: item.GroupShortcut,
			CreatedBy:     getUsername(c),
		}
		
		if err := database.DB.Create(&newGroup).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("CREATE_FAIL")
			return c.JSON(response)
		}
		
		tempNewGroup := tempGroup(newGroup)
		if err := tempNewGroup.AfterCreate(tx); err != nil {
			response.Status = false
			response.Message = config.GetMessageCode("AFTER_CREATE_FAIL")
			return c.JSON(response)
		}
		

	}

	response.Status = true
	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	return c.JSON(response)
}

// UpdateGroup cập nhật thông tin một Group dựa trên ID
// @Summary Cập nhật thông tin Group
// @Description Cập nhật thông tin một Group dựa trên ID
// @Tags Group
// @Accept json
// @Produce json
// @Param body body []model.UpdateGroupModel true "Thông tin Group cần cập nhật"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group [put]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func UpdateGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateGroupModel
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {
		listCheck := []string{"GroupNameVN", "GroupNameEN", "GroupNameJP"}
		vItem := map[string]string{"GroupNameVN": item.GroupNameVN, "GroupNameEN": item.GroupNameEN, "GroupNameJP": item.GroupNameJP}
		errors := utils.RequireCheck(listCheck, vItem, map[string]string{})

		if len(errors) > 0 {
			response.Message = "MISSING_FIELDS" //NOT_ID_EXISTS
			response.ValidateError = errors
			return c.JSON(response)
		}

		tempNGroup := updateGr(*item)
		
		if err := tempNGroup.BeforeUpdate(tx); err != nil {
			response.Status = false
			response.Message = config.GetMessageCode("ERROR_BEFORE_UPDATE")
			return c.JSON(response)
		}

		var group model.Group
		if item.GroupID != 0 {
			if item.IsDeleted {
				// Soft delete group
				results := db.First(&group, item.GroupID)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("NOT_ID_EXISTS")
					return c.JSON(response)
				}
				group.DeletedBy = getUsername(c)
				group.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			} else {
				// Update group
				results := db.First(&group, item.GroupID)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("NOT_ID_EXISTS")
					return c.JSON(response)
				}

				group.DepartmentID = item.DepartmentID
				group.GroupNameVN = item.GroupNameVN
				group.GroupNameEN = item.GroupNameEN
				group.GroupNameJP = item.GroupNameJP
				group.GroupShortcut = item.GroupShortcut


				if err := db.Model(&group).Updates(group).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}

			if err := db.Save(&group).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("SYSTEM_ERROR")
				return c.JSON(response)
			}
		} else {
			// Create group
			newGroup := model.Group{
				GroupNameVN:   item.GroupNameVN,
				GroupNameEN:   item.GroupNameEN,
				GroupNameJP:   item.GroupNameJP,
				GroupShortcut: item.GroupShortcut,
				CreatedBy:     getUsername(c),
			}

			if err := db.Create(&newGroup).Error; err != nil {
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

// DeleteGroup xóa một Group dựa trên ID
// @Summary Xóa Group
// @Description Xóa một Group dựa trên ID
// @Tags Group
// @Accept json
// @Produce json
// @Param id path int true "ID của Group"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group/{id} [delete]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func DeleteGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	groupID := c.Params("id")

	// Check if the group exists
	var group model.Group
	result := database.DB.First(&group, groupID)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}


	results := database.DB.First(&group, "group_id = ?", groupID)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Soft delete data: Update deleted_at field with the current time
	group.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	group.DeletedBy = getUsername(c)

	// Update data: update group by id from group table (database)
	if err := database.DB.Model(&group).Updates(&group).Error; err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}


	response.Status = true
	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	return c.JSON(response)
}

// RestoreGroup khôi phục một Group dựa trên ID
// @Summary Khôi phục Group
// @Description Khôi phục một Group dựa trên ID
// @Tags Group
// @Accept json
// @Produce json
// @Param id path int true "ID của Group"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /group/{id}/restore [put]
// @Security ApiKeyAuth
// @Security ApiTokenAuth
func RestoreGroup(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	groupID := c.Params("id")

	// Check if the group exists
	var group model.Group
	result := database.DB.Unscoped().First(&group, groupID)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}


	// Restore the group
	if err := database.DB.Unscoped().Model(&group).Update("deleted_at", nil).Error; err != nil {
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

	var department modell.Department
	results := database.DB.Select("*").Where("department_id = ?", g.DepartmentID).First(&department)
	if results.Error != nil {
		return results.Error
	}

	if !tx.Statement.Changed("group_id", "group_name_vn", "group_name_en", "group_name_jp", "group_shortcut") {
		WriteLog("DATA_NOT_CHANGED")
	}

	if tx.Statement.Changed("LogVersion") {
		return errors.New(config.Config("DATA_NOT_CHANGED"))
	}

	group := tx.Statement.Dest.(*model.Group)
	group.LogVersion++


	return nil
}



func (g *tempGroup) AfterCreate(tx *gorm.DB) (err error) {
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
