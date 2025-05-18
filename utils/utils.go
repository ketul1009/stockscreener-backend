package utils

// Define a generic Set type
type Set[T comparable] struct {
	data map[T]struct{}
}

// Create a new Set
func NewSet[T comparable](items ...T) *Set[T] {
	set := &Set[T]{data: make(map[T]struct{})}
	for _, item := range items {
		set.Add(item)
	}
	return set
}

// Add an item to the Set
func (s *Set[T]) Add(item T) {
	s.data[item] = struct{}{}
}

// Remove an item from the Set
func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}

// Check if item exists in the Set
func (s *Set[T]) Has(item T) bool {
	_, exists := s.data[item]
	return exists
}

// Return the number of items
func (s *Set[T]) Size() int {
	return len(s.data)
}

// Get all items as a slice
func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.data))
	for item := range s.data {
		items = append(items, item)
	}
	return items
}
