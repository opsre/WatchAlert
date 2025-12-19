package initialization

import (
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

var perms []models.UserPermissions

func InitPermissionsSQL(ctx *ctx.Context) {
	var psData []models.UserPermissions

	for _, v := range models.PermissionsInfo() {
		psData = append(psData, v)
	}
	perms = psData

	ctx.DB.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.UserPermissions{})
	ctx.DB.DB().Model(&models.UserPermissions{}).Create(&psData)
}

func InitUserRolesSQL(ctx *ctx.Context) {
	var adminRole models.UserRole
	var db = ctx.DB.DB()

	roles := models.UserRole{
		ID:          "admin",
		Name:        "admin",
		Description: "system",
		Permissions: perms,
		UpdateAt:    time.Now().Unix(),
	}

	err := db.Model(&models.UserRole{}).Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&roles).Error
		}
	} else {
		err = db.Where("name = ?", "admin").Updates(models.UserRole{Permissions: perms}).Error
	}

	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		panic(err)
	}
}
