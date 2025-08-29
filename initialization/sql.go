package initialization

import (
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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
	var db = ctx.DB.DB().Model(&models.UserRole{})

	roles := models.UserRole{
		ID:          "admin",
		Name:        "admin",
		Description: "system",
		Permissions: perms,
		CreateAt:    time.Now().Unix(),
	}

	err := db.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ctx.DB.DB().Create(&roles).Error
		}
	} else {
		err = db.Where("name = ?", "admin").Updates(models.UserRole{Permissions: perms}).Error
	}

	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		panic(err)
	}
}
