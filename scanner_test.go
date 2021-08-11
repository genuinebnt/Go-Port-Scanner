package scanner

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	got := Scanner()
	want := []int{22, 80}
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
