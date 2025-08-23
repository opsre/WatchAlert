package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	NoticeTmplRepo struct {
		entryRepo
	}

	InterNoticeTmplRepo interface {
		List(id, noticeType, query string) ([]models.NoticeTemplateExample, error)
		Create(r models.NoticeTemplateExample) error
		Update(r models.NoticeTemplateExample) error
		Delete(id string) error
		Get(id string) models.NoticeTemplateExample
	}
)

func newNoticeTmplInterface(db *gorm.DB, g InterGormDBCli) InterNoticeTmplRepo {
	return &NoticeTmplRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (nr NoticeTmplRepo) List(id, noticeType, query string) ([]models.NoticeTemplateExample, error) {
	var (
		data []models.NoticeTemplateExample
		db   = nr.db.Model(&models.NoticeTemplateExample{})
	)

	if id != "" {
		db.Where("id = ?", id)
	}
	if noticeType != "" {
		db.Where("notice_type = ?", noticeType)
	}
	if query != "" {
		db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (nr NoticeTmplRepo) Create(r models.NoticeTemplateExample) error {
	err := nr.g.Create(models.NoticeTemplateExample{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Update(r models.NoticeTemplateExample) error {
	u := Updates{
		Table: models.NoticeTemplateExample{},
		Where: map[string]interface{}{
			"id = ?": r.ID,
		},
		Updates: r,
	}

	err := nr.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Delete(id string) error {
	d := Delete{
		Table: models.NoticeTemplateExample{},
		Where: map[string]interface{}{
			"id = ?": id,
		},
	}

	err := nr.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Get(id string) models.NoticeTemplateExample {
	var (
		data models.NoticeTemplateExample
		db   = nr.db.Model(&models.NoticeTemplateExample{})
	)
	if id != "" {
		db.Where("id = ?", id)
	}

	err := db.First(&data).Error
	if err != nil {
		return data
	}
	return data
}
