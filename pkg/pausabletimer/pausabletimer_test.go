package pausabletimer

import (
	"testing"
	"time"
)

type test_state struct {
	S bool
}

func newTestState() *test_state {
	return &test_state{false}
}
func (s *test_state) setTrue() {
	s.S = true
}
func (s *test_state) setFalse() {
	s.S = false
}

func TestPauseResumeZero(t *testing.T) {
	state := newTestState()

	pausable := NewPausableTimer(100, state.setTrue)
	if !pausable.Pause() {
		t.Error("Unexpected error Pause() failed")
		return
	}

	testtimer := time.NewTimer(time.Duration(msToNs(100)))
	<-testtimer.C

	pausable.Resume()

	if !state.S {
		t.Error("Expected state.S to be true, got false")
	}
}
func TestPauseResumeNegative(t *testing.T) {
	state := newTestState()

	pausable := NewPausableTimer(100, state.setTrue)
	if !pausable.Pause() {
		t.Error("Unexpected error Pause() failed")
		return
	}

	testtimer := time.NewTimer(time.Duration(msToNs(1000)))
	<-testtimer.C

	pausable.Resume()

	if !state.S {
		t.Error("Expected state.S to be true, got false")
	}
}

func TestPauseResumeMutation(t *testing.T) {
	state := newTestState()

	pausable := NewPausableTimer(100, state.setTrue)
	if !pausable.Pause() {
		t.Error("Unexpected error Pause() failed")
		return
	}

	testtimer := time.NewTimer(time.Duration(msToNs(100)))
	<-testtimer.C

	pausable.Resume()
	state.setFalse()
	pausable.Pause()
	pausable.Resume()

	if state.S {
		t.Error("Expected state.S to be false, unmutated, got true")
	}
}

func TestResetWhilePaused(t *testing.T) {
	state := newTestState()

	pausable := NewPausableTimer(100, state.setTrue)
	if !pausable.Pause() {
		t.Error("Unexpected error Pause() failed")
		return
	}

	testtimer := time.NewTimer(time.Duration(msToNs(100)))
	<-testtimer.C

	pausable.Reset(100)

	// Wait a bit
	testtimer = time.NewTimer(time.Duration(msToNs(10)))
	<-testtimer.C

	if state.S {
		t.Error("Expected state.S to be false, unmutated, got true")
	}
}
