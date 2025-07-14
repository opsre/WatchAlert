package repo

import (
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	CommentRepo struct {
		entryRepo
	}

	InterCommentRepo interface {
		Add(r models.RequestAddEventComment) error
		Delete(r models.RequestDeleteEventComment) error
		List(r models.RequestListEventComments) ([]models.Comment, error)
	}
)

func newCommentInterface(db *gorm.DB, g InterGormDBCli) InterCommentRepo {
	return &CommentRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (c CommentRepo) List(r models.RequestListEventComments) ([]models.Comment, error) {
	var data = []models.Comment{}

	db := c.db.Model(&models.Comment{})
	db.Where("tenant_id = ? AND fingerprint = ?", r.TenantId, r.Fingerprint)
	if err := db.Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (c CommentRepo) Add(r models.RequestAddEventComment) error {
	db := c.db.Model(&models.Comment{})

	return db.Create(&models.Comment{
		TenantId:    r.TenantId,
		CommentId:   tools.RandUid(),
		Fingerprint: r.Fingerprint,
		Username:    r.Username,
		UserId:      r.UserId,
		Time:        time.Now().Unix(),
		Content:     r.Content,
	}).Error
}

func (c CommentRepo) Delete(r models.RequestDeleteEventComment) error {
	db := c.db.Model(&models.Comment{})
	return db.Where("tenant_id = ? AND comment_id = ?", r.TenantId, r.CommentId).Delete(&models.Comment{}).Error
}
