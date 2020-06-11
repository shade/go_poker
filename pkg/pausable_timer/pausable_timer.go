package pausable_timer

import (
	"time"
)

// PausableTimer describes a timer that can be paused
type PausableTimer struct {
	ending time.Time

	timer  *time.Timer
	paused bool
	cb     func()
}

// Converts milliseconds to nanoseconds
func msToNs(ms int64) int64 {
	return ms * 1000000
}

// NewPausableTimer creates a PausableTimer based on a millisecond
// count and a callback set to run after
func NewPausableTimer(ms int64, cb func()) *PausableTimer {
	t := time.Duration(msToNs(ms))

	now := time.Now()
	timer := time.AfterFunc(t, cb)

	return &PausableTimer{
		now.Add(t),
		timer,
		false,
		cb,
	}
}

// Reset resets the timer to call the existing callback in an additional ms milliseconds
// This method is only callable on paused or expired timers
func (p *PausableTimer) Reset(ms int64) bool {
	if !p.isExpired() && !p.paused {
		return false
	}
	t := time.Duration(msToNs(ms))

	p.paused = false
	p.ending = time.Now().Add(t)
	p.timer.Reset(t)

	return true
}

// Pause pauses the timer until Resume is called again.
// Cannot be called on paused or expired timers
func (p *PausableTimer) Pause() bool {
	if p.paused || p.isExpired() {
		return false
	}

	if !p.timer.Stop() {
		return false
	}

	p.paused = true
	return true
}

// Resume is called to resume a paused timer ONLY.
// Expired timers call the callback immediately on stack.
func (p *PausableTimer) Resume() bool {
	if !p.paused {
		return false
	}
	remainder := time.Until(p.ending)

	// Execute non-expired callbacks immediately
	if p.isExpired() {
		p.cb()
	} else {
		p.timer.Reset(remainder)
	}
	p.paused = false

	return true
}

func (p *PausableTimer) isExpired() bool {
	return time.Until(p.ending) <= 0
}
