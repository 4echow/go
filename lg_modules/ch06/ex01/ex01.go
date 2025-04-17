package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age       int
}

func MakePerson(firstName, lastName string, age int) Person {
	person := Person{
		firstName,
		lastName,
		age,
	}
	return person
}

func MakePersonPointer(firstName, lastName string, age int) *Person {
	person := Person{
		firstName,
		lastName,
		age,
	}
	return &person
}

func main() {
	myPerson := Person{
		firstName: "Bob",
		lastName:  "Baumeister",
		age:       42,
	}

	prs1 := MakePerson(myPerson.firstName, myPerson.lastName, myPerson.age)
	fmt.Println(prs1)
	prs2 := MakePersonPointer(myPerson.firstName, myPerson.lastName, myPerson.age)
	fmt.Println(prs2)
}
