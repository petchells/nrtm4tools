package util

import (
	"fmt"
	"sort"
	"testing"
)

func TestStringSet(t *testing.T) {
	s := NewSet[string]()
	if !s.IsEmpty() {
		t.Error("Expected an empty set but was", s.String())
	}
	if s.String() != "[]" {
		t.Error("Expected a string showing an empty set but was", s.String())
	}
	s.Add("A")
	s.Add("B")
	expected := "[A B]"
	result := s.Members()
	sort.Strings(result)
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
	sort.Strings(result)
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

func TestSetUnion(t *testing.T) {
	s1 := NewSet("B", "C", "A")
	s2 := NewSet("B", "D", "A")
	result := s1.Union(s2).Members()
	sort.Strings(result)
	expected := "[A B C D]"
	if expected != fmt.Sprint(result) {
		t.Error("Expected", expected, "but got", result)
	}
}

func TestSetContains(t *testing.T) {
	s1 := NewSet("B", "C", "A")
	s2 := NewSet("B", "D", "A")
	s3 := NewSet("B", "A")
	{
		result := s1.ContainsValues(s2.Members())
		expected := false
		if expected != result {
			t.Error("Expected", expected, "but got", result)
		}
	}
	{
		result := s1.ContainsValues(s3.Members())
		expected := true
		if expected != result {
			t.Error("Expected", expected, "but got", result)
		}
	}
}
