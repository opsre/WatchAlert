package services

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type eventService struct {
	ctx *ctx.Context
}

type InterEventService interface {
	ListCurrentEvent(req interface{}) (interface{}, interface{})
	ListHistoryEvent(req interface{}) (interface{}, interface{})
	ProcessAlertEvent(req interface{}) (interface{}, interface{})
	ListComments(req interface{}) (interface{}, interface{})
	AddComment(req interface{}) (interface{}, interface{})
	DeleteComment(req interface{}) (interface{}, interface{})
}

func newInterEventService(ctx *ctx.Context) InterEventService {
	return &eventService{
		ctx: ctx,
	}
}

func (e eventService) ProcessAlertEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestProcessAlertEvent)

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
	r, ok := req.(*types.RequestAlertCurEventQuery)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *models.AlertCurEventQuery")
	}

	center, err := e.ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(r.TenantId, r.FaultCenterId))
	if err != nil {
		return nil, err
	}

	var (
		allEvents      []models.AlertCurEvent
		filteredEvents []models.AlertCurEvent
		curTime        = time.Now()
	)
	for _, alert := range center {
		allEvents = append(allEvents, *alert)
	}

	var form int64
	var to int64
	if r.Scope > 0 {
		to = curTime.Unix()
		form = curTime.Add(-time.Duration(r.Scope) * 24 * time.Hour).Unix()
	}

	for _, event := range allEvents {
		if r.DatasourceType != "" && event.DatasourceType != r.DatasourceType {
			continue
		}

		if r.Severity != "" && event.Severity != r.Severity {
			continue
		}

		if r.Scope > 0 && (event.FirstTriggerTime < form || event.FirstTriggerTime > to) {
			continue
		}

		if r.Query != "" {
			queryMatch := false
			if strings.Contains(event.RuleName, r.Query) {
				queryMatch = true
			} else if strings.Contains(event.Annotations, r.Query) {
				queryMatch = true
			} else if event.Labels != nil && strings.Contains(tools.JsonMarshalToString(event.Labels), r.Query) {
				queryMatch = true
			}
			if !queryMatch {
				continue
			}
		}

		if r.FaultCenterId != "" && !strings.Contains(event.FaultCenterId, r.FaultCenterId) {
			continue
		}

		if r.Status != "" && string(event.Status) != r.Status {
			continue
		}

		filteredEvents = append(filteredEvents, event)
	}

	sort.Slice(filteredEvents, func(i, j int) bool {
		a, b := &filteredEvents[i], &filteredEvents[j]

		// 按持续时间降序
		durA := a.LastEvalTime - a.FirstTriggerTime
		durB := b.LastEvalTime - b.FirstTriggerTime
		switch r.SortOrder {
		case models.SortOrderASC:
			if durA != durB {
				return durA < durB // 升序
			}
		case models.SortOrderDesc:
			if durA != durB {
				return durA > durB // 降序
			}
		default:
			if a.FirstTriggerTime != b.FirstTriggerTime {
				return a.Fingerprint < b.Fingerprint
			}
		}

		// 默认按指纹升序
		return a.Fingerprint < b.Fingerprint
	})

	paginatedList := pageSlice(filteredEvents, int(r.Page.Index), int(r.Page.Size))
	return types.ResponseAlertCurEventList{
		List: paginatedList,
		Page: models.Page{
			Total: int64(len(filteredEvents)),
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (e eventService) ListHistoryEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertHisEventQuery)
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

func (e eventService) ListComments(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestListEventComments)
	comment := e.ctx.DB.Comment()
	data, err := comment.List(*r)
	if err != nil {
		return nil, fmt.Errorf("获取评论失败, %s", err.Error())
	}

	return data, nil
}

func (e eventService) AddComment(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAddEventComment)
	comment := e.ctx.DB.Comment()
	err := comment.Add(*r)
	if err != nil {
		return nil, fmt.Errorf("评论失败, %s", err.Error())
	}

	return "评论成功", nil
}

func (e eventService) DeleteComment(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDeleteEventComment)
	comment := e.ctx.DB.Comment()
	err := comment.Delete(*r)
	if err != nil {
		return nil, fmt.Errorf("删除评论失败, %s", err.Error())
	}

	return "删除评论成功", nil
}
