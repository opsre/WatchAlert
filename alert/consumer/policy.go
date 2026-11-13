package consumer

import (
	"sync"
	"time"
)

type consumerRestartState struct {
	attempts  int
	lastPanic time.Time
}

type consumerRestartPolicy struct {
	mu     sync.Mutex
	states map[string]consumerRestartState
}

func (p *consumerRestartPolicy) next(faultCenterID string, now time.Time) (time.Duration, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.states == nil {
		p.states = make(map[string]consumerRestartState)
	}

	state := p.states[faultCenterID]
	if !state.lastPanic.IsZero() && now.Sub(state.lastPanic) >= ConsumerRestartResetAfter {
		state.attempts = 0
	}
	if state.attempts >= ConsumerRestartMaxAttempts {
		return 0, false
	}

	state.attempts++
	state.lastPanic = now
	p.states[faultCenterID] = state

	delay := ConsumerRestartBaseDelay
	for i := 1; i < state.attempts; i++ {
		delay *= 2
		if delay >= ConsumerRestartMaxDelay {
			return ConsumerRestartMaxDelay, true
		}
	}
	return delay, true
}
