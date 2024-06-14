package main

import "fmt"

type iParent interface {
	//ProcessSomething(st iParent) int
	//GetID() int
}

type Parent struct {
	ID int
}

func (p *Parent) ProcessSomething(st iParent) int {
	return p.ID
}

func (p *Parent) SetID(id int) {
	p.ID = id
}

type Child struct {
	Parent
	Name string
}

/* func (c *Child) GetID() int {
	return 0
} */

/*
	 func (c *Child) ProcessSomething(st iParent) int {
		return 0
	}
*/
func main() {

	parent := &Parent{}
	child := &Child{}
	child.SetID(2)
	id := parent.ProcessSomething(child)
	fmt.Printf("id: %d\n", id)
	fmt.Printf("parent id: %d\n", parent.ID)
	fmt.Printf("child id: %d\n", child.ID)
}
