package set

// Define a generic Set type
type Set[T comparable] map[T]struct{}

// Add an element to the set
func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

// Remove an element from the set
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Check if an element is in the set
func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

// Get the number of elements in the set
func (s Set[T]) Size() int {
	return len(s)
}

// Convert the set to a slice
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for element := range s {
		result = append(result, element)
	}
	return result
}

func ConvertArray[T comparable](myArray []T) Set[T] {
	mySet := make(Set[T])
	for _, element := range myArray {
		mySet.Add(element)
	}
	return mySet
}
