package util

import (
	"fmt"
	"sort"
	"testing"
)

func TestStringSet(t *testing.T) {
	s := NewSet[string]()
	s.Add("A")
	s.Add("B")
	expected := "[A B]"
	result := s.Members()
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	resultStr := fmt.Sprintf("%v", result)
	if expected != resultStr {
		t.Error("Expected", expected, "but got", resultStr)
	}
}

func TestSetIntersection(t *testing.T) {
	s1 := NewSet("A", "B")
	s2 := NewSet("B", "C")
	expected := "[B]"
	result := s1.Intersection(s2).Members()
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	resultStr := fmt.Sprintf("%v", result)
	if expected != resultStr {
		t.Error("Expected", expected, "but got", resultStr)
	}
}

func TestSetDifference(t *testing.T) {
	s1 := NewSet("A", "B")
	s2 := NewSet("B", "C", "A")
	expected := "[]"
	result := s1.Difference(s2)
	resultStr := fmt.Sprintf("%v", result.Members())
	if expected != resultStr {
		t.Error("Expected", expected, "but got", resultStr)
	}
}

func TestSetFilter(t *testing.T) {
	s1 := NewSet("B", "C", "A")
	result := s1.Filter(func(ele string) bool {
		return ele != "C"
	}).Members()
	sort.Strings(result)
	expected := "[A B]"
	if expected != fmt.Sprint(result) {
		t.Error("Expected", expected, "but got", result)
	}
}
