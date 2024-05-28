package util

import "fmt"

type Set[E comparable] map[E]struct{}
type FilterFn[E comparable] func(E) bool

//type MapFn[E comparable, F comparable] func(E) F

func NewSet[E comparable](vals ...E) Set[E] {
	s := Set[E]{}
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[E]) Add(vals ...E) {
	for _, v := range vals {
		s[v] = struct{}{}
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

func (s Set[E]) Filter(fn FilterFn[E]) Set[E] {
	result := NewSet[E]()
	for _, ele := range s.Members() {
		if fn(ele) {
			result.Add(ele)
		}
	}
	return result
}
