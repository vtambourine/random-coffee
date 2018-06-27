package tests

import "testing"

func TestMatch(t *testing.T) {
	employees := []Employee{
		Employee{
			name: "Alice",
			slots: []TimeSlot{9, 10, 11},
		},
		Employee{
			name: "Bob",
			slots: []TimeSlot{10, 11, 12},
		},
	}

	GetMatch(employees)

}

func TestNonMatch(t *testing.T) {
	t.Error("wow wo")
}