package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	EventRepo struct {
		entryRepo
	}

	InterEventRepo interface {
		GetHistoryEvent(r models.AlertHisEventQuery) (models.HistoryEventResponse, error)
		CreateHistoryEvent(r models.AlertHisEvent) error
	}
)

func newEventInterface(db *gorm.DB, g InterGormDBCli) InterEventRepo {
	return &EventRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (e EventRepo) GetHistoryEvent(r models.AlertHisEventQuery) (models.HistoryEventResponse, error) {
	var data []models.AlertHisEvent
	var count int64

	db := e.DB().Model(&models.AlertHisEvent{})
	db.Where("tenant_id = ?", r.TenantId)
	db.Where("fault_center_id = ?", r.FaultCenterId)

	if r.Query != "" {
		db.Where("rule_name LIKE ? OR severity LIKE ? OR annotations LIKE ? OR fingerprint LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}

	if r.DatasourceType != "" {
		db = db.Where("datasource_type = ?", r.DatasourceType)
	}

	if r.Severity != "" {
		db = db.Where("severity = ?", r.Severity)
	}

	if r.StartAt != 0 && r.EndAt != 0 {
		db = db.Where("first_trigger_time > ? and first_trigger_time < ?", r.StartAt, r.EndAt)
	}

	if err := db.Count(&count).Error; err != nil {
		return models.HistoryEventResponse{}, err
	}

	switch r.SortOrder {
	case models.SortOrderASC:
		db.Order("alarm_duration asc")
	case models.SortOrderDesc:
		db.Order("alarm_duration desc")
	default:
		db.Order("recover_time desc")
	}

	if err := db.Limit(int(r.Page.Size)).Offset(int((r.Page.Index - 1) * r.Page.Size)).Find(&data).Error; err != nil {
		return models.HistoryEventResponse{}, err
	}

	return models.HistoryEventResponse{
		List: data,
		Page: models.Page{
			Index: r.Page.Index,
			Size:  r.Page.Size,
			Total: count,
		},
	}, nil
}

func (e EventRepo) CreateHistoryEvent(r models.AlertHisEvent) error {
	err := e.g.Create(models.AlertHisEvent{}, r)
	if err != nil {
		return err
	}

	return nil
}
