package storage

import (
	"testing"
)

//TestNewStorage ...
func TestNewStorage(t *testing.T) {
	s := NewStorage()
	if s == nil {
		t.Fail()
	}
}
