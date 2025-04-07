package main

import "fmt"

type Animal interface {
	Speak() string
	Run() string
}

func newAnimal(a Animal) {
	fmt.Println(a.Speak())
	fmt.Println(a.Run())
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}
func (d Dog) Run() string {
	return "running"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Run() string {
	return "running"
}

func main() {
	dog := Dog{Name: "Buddy"}
	cat := Cat{Name: "Whiskers"}

	newAnimal(dog)
	newAnimal(cat)
}
