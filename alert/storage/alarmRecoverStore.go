package storage

import (
	"sort"
	"sync"
)

type AlarmRecoverWaitStore struct {
	data  []RuleEntry
	mutex sync.RWMutex
}

type RuleEntry struct {
	RuleID       string
	Fingerprints map[string]int64
}

func NewAlarmRecoverStore() *AlarmRecoverWaitStore {
	return &AlarmRecoverWaitStore{
		data: make([]RuleEntry, 0),
	}
}

func (a *AlarmRecoverWaitStore) findRuleEntryPos(ruleID string) int {
	return sort.Search(len(a.data), func(i int) bool {
		return a.data[i].RuleID >= ruleID
	})
}

func (a *AlarmRecoverWaitStore) Set(ruleID string, fingerprint string, t int64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pos := a.findRuleEntryPos(ruleID)

	// 更新现有条目
	if pos < len(a.data) && a.data[pos].RuleID == ruleID {
		if a.data[pos].Fingerprints == nil {
			a.data[pos].Fingerprints = make(map[string]int64)
		}
		a.data[pos].Fingerprints[fingerprint] = t
		return
	}

	// 插入新条目（保持有序）
	a.data = append(a.data, RuleEntry{})
	copy(a.data[pos+1:], a.data[pos:])
	a.data[pos] = RuleEntry{
		RuleID:       ruleID,
		Fingerprints: map[string]int64{fingerprint: t},
	}
}

func (a *AlarmRecoverWaitStore) Get(ruleID string, fingerprint string) (int64, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	pos := a.findRuleEntryPos(ruleID)
	if pos < len(a.data) && a.data[pos].RuleID == ruleID {
		val, exists := a.data[pos].Fingerprints[fingerprint]
		return val, exists
	}
	return 0, false
}

func (a *AlarmRecoverWaitStore) Remove(ruleID string, fingerprint string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pos := a.findRuleEntryPos(ruleID)
	if pos < len(a.data) && a.data[pos].RuleID == ruleID {
		delete(a.data[pos].Fingerprints, fingerprint)
	}
}

func (a *AlarmRecoverWaitStore) List(ruleID string) map[string]int64 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	pos := a.findRuleEntryPos(ruleID)
	if pos < len(a.data) && a.data[pos].RuleID == ruleID {
		copyMap := make(map[string]int64, len(a.data[pos].Fingerprints))
		for k, v := range a.data[pos].Fingerprints {
			copyMap[k] = v
		}
		return copyMap
	}
	return make(map[string]int64)
}
