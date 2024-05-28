package util

import (
	"fmt"
	"testing"
)

type MockIntItem struct {
	ID   int
	Name string
}

func TestHashMapper(t *testing.T) {
	list := []MockIntItem{
		{ID: 5, Name: "Five"},
		{ID: 6, Name: "Six"},
	}
	expected := "map[5:{5 Five} 6:{6 Six}]"
	result := SliceToMap(func(e MockIntItem) int {
		return e.ID
	}, list)
	resultStr := fmt.Sprint(result)
	if expected != resultStr {
		t.Errorf("Expected %v but got %v", expected, resultStr)
	}
	t.Log("OK")
}
