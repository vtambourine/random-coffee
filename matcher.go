package main

import (
	"fmt"
)

type Graph struct {
	queue []Employee
}

func NewMatcher() *Graph {
	return &Graph{}
}

func SortInput() {
	fmt.Printf("Tak")
}

func (m *Graph) Add(e *Employee) {
	//m.queue = append(m.queue, e)
}

func (m *Graph) Remove(e *Employee) {
	// Remove employee
}

func (m *Graph) GetMatches() [][]Employee {
	groups := make(map[Office][]Employee)

	for _, e := range m.queue {
		//fmt.Printf("%s at %s\n", e.Name, e.Office.String())
		groups[e.Office] = append(groups[e.Office], e)

	}

	for o, e := range groups {
		fmt.Printf("%s = ", o)
		for _, i := range e {
			fmt.Printf("%v ", i.Name)
		}
		fmt.Printf("\n")
	}

	//fmt.Printf("len %#v\n", m.queue)
	return make([][]Employee, 4)
}
