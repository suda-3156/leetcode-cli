package set

// Set is a generic set data structure.
type Set[T comparable] map[T]struct{}

// New creates a new set with the given values.
func New[T comparable](vals ...T) Set[T] {
	s := make(Set[T])
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

// Add adds an element to the set.
func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

// Remove removes an element from the set.
// If the element does not exist, it does nothing.
func (s *Set[T]) Remove(v T) {
	delete(*s, v)
}

// Contains checks if the set contains the given element.
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// ToSlice converts the set to a slice.
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for v := range s {
		result = append(result, v)
	}
	return result
}

// Union(s1, s2) = s1 ∪ s2
func Union[T comparable](s1, s2 Set[T]) Set[T] {
	result := New[T]()
	for v := range s1 {
		result.Add(v)
	}
	for v := range s2 {
		result.Add(v)
	}
	return result
}

// Intersection(s1, s2) = s1 ∩ s2
func Intersection[T comparable](s1, s2 Set[T]) Set[T] {
	result := New[T]()
	for v := range s1 {
		if s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference(s1, s2) = s1 \ s2
func Difference[T comparable](s1, s2 Set[T]) Set[T] {
	result := New[T]()
	for v := range s1 {
		if !s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
