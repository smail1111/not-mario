package utils

type Set struct {
	items map[any]struct{}
}

// Adds an item to a Set.
func (s *Set) Add(item any) {
	s.items[item] = struct{}{}
}

// Removes an item from a Set.
func (s *Set) Remove(item any) {
	delete(s.items, item)
}

// Removes every item in a Set.
func (s *Set) Clear() {
	s.items = map[any]struct{}{}
}

// Returns whether an item is in a Set.
func (s Set) Has(item any) (has bool) {
	_, has = s.items[item]

	return
}

// Returns an unordered slice of all of the items in a Set.
func (s Set) Items() (items []any) {
	for item := range s.items {
		items = append(items, item)
	}

	return
}

// Returns a new Set containing the given items.
func NewSet(items ...any) (set Set) {
	set.items = map[any]struct{}{}

	for _, item := range items {
		set.Add(item)
	}

	return
}
