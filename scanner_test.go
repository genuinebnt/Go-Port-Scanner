package scanner

import (
	"reflect"
	"testing"
)

func validate(t testing.TB, ports, want []int,) {
	if !reflect.DeepEqual(ports, want) {
		t.Errorf("got %v want %v", ports, want)
	}
}

func TestScanner(t *testing.T) {
	t.Run("test with nmap url and variable ports", func(t *testing.T) {
		url := "scanme.nmap.org"
		got := Scanner(url, "33, 100, 45-47, 22 - 23, 80")
		want := []int{22, 80}
		
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestCheckValidPort(t *testing.T) {
	var ports []int
	_ = checkValidPort("65535", &ports)
	want := []int{65535}
	
	validate(t, ports, want)
}

func TestDashSplit(t *testing.T) {
	var ports []int
	_ = dashSplit("11-13", &ports)
	want := []int{11, 12, 13}
	
	validate(t, ports, want)
}

func TestSplitString(t *testing.T) {
	t.Run("testing with only commma", func(t *testing.T) {
		var ports []int
		ports, _ = stringSplit("12, 14")
		want := []int{12, 14}
		
		validate(t, ports, want)
	})
	
	t.Run("testing with dashes", func(t *testing.T){
		var ports []int
		ports, _ = stringSplit("33 - 35")
		want := []int{33, 34, 35}
		
		validate(t, ports, want)
		
	})
	
	t.Run("testing with comma and dashes", func(t *testing.T){
		var ports []int
		ports, _ = stringSplit("33, 100, 45-47, 22 - 23")
		want := []int{33, 100, 45, 46, 47, 22, 23}
		
		validate(t, ports,want)
	})
}