// Package user provides user domain types and a hardcoded in-memory store.
package user

// User represents a person in the system.
// Value semantics are used throughout: User is small and does not share mutable state.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Store holds the set of known users.
// Zero value is not useful here; construct via New.
type Store struct {
	users []User
}

// New constructs a Store pre-populated with a hardcoded list.
// Dependencies are passed explicitly; no global state is used.
func New() *Store {
	return &Store{
		users: []User{
			{ID: "1", Name: "Felix García",     Role: "engineer"},
			{ID: "2", Name: "Luis Martínez",    Role: "engineer"},
			{ID: "3", Name: "Alejandro Reyes",  Role: "engineer"},
			{ID: "4", Name: "Robizon Kharebava", Role: "engineer"},
			{ID: "5", Name: "Julio Pérez",       Role: "lead"},
		},
	}
}

// All returns every user in the store.
// The returned slice is a copy; callers cannot mutate internal state.
func (s *Store) All() []User {
	out := make([]User, len(s.users))
	copy(out, s.users)
	return out
}
