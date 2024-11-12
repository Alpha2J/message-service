package helper

import "testing"

func TestAdd(t *testing.T) {
	s := NewSet()
	s.Add(1)
	if !s.Contains(1) {
		t.Error("Add failed")
	}
}

func TestRemove(t *testing.T) {
	s := NewSet()
	s.Add(1)
	s.Remove(1)
	if s.Contains(1) {
		t.Error("Remove failed")
	}
}

func TestSize(t *testing.T) {
	s := NewSet()
	s.Add(1)
	if s.Size() != 1 {
		t.Error("Size failed")
	}
}
