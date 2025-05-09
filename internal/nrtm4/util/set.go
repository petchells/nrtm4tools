package util

import "fmt"

type member struct{}

// Set is a mathematical set of values
type Set[E comparable] map[E]member

type filterFn[E comparable] func(E) bool

// NewSet creates a new set from a list of values, or an empty set of a certain type
func NewSet[E comparable](vals ...E) Set[E] {
	s := make(Set[E], len(vals))
	for _, v := range vals {
		s[v] = member{}
	}
	return s
}

func (s Set[E]) Add(vals ...E) {
	for _, v := range vals {
		s[v] = member{}
	}
}

func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

func (s Set[E]) ContainsValues(vals []E) bool {
	for _, v := range vals {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

func (s Set[E]) Members() []E {
	result := make([]E, 0, len(s))
	for v := range s {
		result = append(result, v)
	}
	return result
}

func (s Set[E]) String() string {
	return fmt.Sprintf("%v", s.Members())
}

func (s Set[E]) Union(s2 Set[E]) Set[E] {
	result := NewSet(s.Members()...)
	result.Add(s2.Members()...)
	return result
}

func (s Set[E]) Intersection(s2 Set[E]) Set[E] {
	result := NewSet(s.Members()...)
	for _, v := range s.Members() {
		if !s2.Contains(v) {
			delete(result, v)
		}
	}
	return result
}

func (s Set[E]) Difference(s2 Set[E]) Set[E] {
	result := NewSet(s.Members()...)
	for _, v := range s.Members() {
		if s2.Contains(v) {
			delete(result, v)
		}
	}
	return result
}

func (s Set[E]) IsEmpty() bool {
	return len(s) == 0
}

func (s Set[E]) Filter(fn filterFn[E]) Set[E] {
	result := NewSet[E]()
	for _, member := range s.Members() {
		if fn(member) {
			result.Add(member)
		}
	}
	return result
}
