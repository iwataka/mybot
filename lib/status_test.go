package mybot

import (
	"testing"
)

func TestCheckTwitterListenDMStatus(t *testing.T) {
	s := NewStatus()
	for i := 0; i < 20; i++ {
		go func() {
			s.LockListenDMRoutine()
			s.UnlockListenDMRoutine()
		}()
	}
	if s.CheckTwitterListenDMStatus() {
		t.Fatalf("%v expected but %v found", false, true)
	}
}

func TestCheckTwitterListenUsersStatus(t *testing.T) {
	s := NewStatus()
	for i := 0; i < 20; i++ {
		go func() {
			s.LockListenUsersRoutine()
			s.UnlockListenUsersRoutine()
		}()
	}
	if s.CheckTwitterListenUsersStatus() {
		t.Fatalf("%v expected but %v found", false, true)
	}
}
