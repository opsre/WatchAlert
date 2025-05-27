package services

import (
	"strings"
	"sync"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type eventService struct {
	ctx *ctx.Context
}

type InterEventService interface {
	ListCurrentEvent(req interface{}) (interface{}, interface{})
	ListHistoryEvent(req interface{}) (interface{}, interface{})
	ProcessAlertEvent(req interface{}) (interface{}, interface{})
}

func newInterEventService(ctx *ctx.Context) InterEventService {
	return &eventService{
		ctx: ctx,
	}
}

func (e eventService) ProcessAlertEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProcessAlertEvent)

	var wg sync.WaitGroup
	wg.Add(len(r.Fingerprints))
	for _, fingerprint := range r.Fingerprints {
		go func(fingerprint string) {
			defer wg.Done()
			cache, err := e.ctx.Redis.Alert().GetEventFromCache(r.TenantId, r.FaultCenterId, fingerprint)
			if err != nil {
				return
			}

			switch r.State {
			case 1:
				if cache.UpgradeState.IsConfirm {
					return
				}

				cache.UpgradeState.IsConfirm = true
				cache.UpgradeState.WhoAreConfirm = r.Username
				cache.UpgradeState.ConfirmOkTime = r.Time
			case 2:
				if !cache.UpgradeState.IsConfirm && cache.UpgradeState.IsHandle {
					return
				}

				cache.UpgradeState.IsHandle = true
				cache.UpgradeState.WhoAreHandle = r.Username
				cache.UpgradeState.HandleOkTime = r.Time
			}

			e.ctx.Redis.Alert().PushAlertEvent(&cache)
		}(fingerprint)
	}

	wg.Wait()
	return nil, nil
}

func (e eventService) ListCurrentEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertCurEventQuery)
	center, err := e.ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(r.TenantId, r.FaultCenterId))
	if err != nil {
		return nil, err
	}

	var dataList []models.AlertCurEvent
	for _, alert := range center {
		dataList = append(dataList, *alert)
	}

	if r.DatasourceType != "" {
		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if v.DatasourceType == r.DatasourceType {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
		}
		dataList = dsTypeDataList
	}

	if r.Severity != "" {
		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if v.Severity == r.Severity {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
		}
		dataList = dsTypeDataList
	}

	if r.Scope > 0 {
		curTime := time.Now()
		to := curTime.Unix()
		form := curTime.Add(-time.Duration(r.Scope) * (time.Hour * 24)).Unix()

		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if v.FirstTriggerTime > form && v.FirstTriggerTime < to {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
		}
		dataList = dsTypeDataList
	}

	if r.Query != "" {
		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if strings.Contains(v.RuleName, r.Query) {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
			if strings.Contains(v.Annotations, r.Query) {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
			if strings.Contains(tools.JsonMarshal(v.Labels), r.Query) {
				dsTypeDataList = append(dsTypeDataList, v)
				continue
			}
		}
		dataList = dsTypeDataList
	}

	if r.FaultCenterId != "" {
		var data []models.AlertCurEvent
		for _, v := range dataList {
			if strings.Contains(v.FaultCenterId, r.FaultCenterId) {
				data = append(data, v)
				continue
			}
		}
		dataList = data
	}

	return models.CurEventResponse{
		List: pageSlice(dataList, int(r.Page.Index), int(r.Page.Size)),
		Page: models.Page{
			Total: int64(len(dataList)),
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil

}

func (e eventService) ListHistoryEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertHisEventQuery)
	data, err := e.ctx.DB.Event().GetHistoryEvent(*r)
	if err != nil {
		return nil, err
	}

	return data, err

}

func pageSlice(data []models.AlertCurEvent, index, size int) []models.AlertCurEvent {
	if index <= 0 {
		index = 1
	}

	if size <= 0 {
		index = 10
	}

	total := len(data)
	if total == 0 {
		return []models.AlertCurEvent{}
	}

	offset := (index - 1) * size
	if offset >= total {
		return []models.AlertCurEvent{}
	}

	limit := index * size
	if limit > total {
		limit = total
	}

	return data[offset:limit]
}
