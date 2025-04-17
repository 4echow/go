package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age       int
}

func MakeManyPersons(person Person) []Person {
	personSlice := make([]Person, 10_000_000)
	for i := range len(personSlice) {
		personSlice[i] = person
	}
	return personSlice
}

func MakeManyPersonsCap(person Person) []Person {
	personSlice := make([]Person, 0, 10_000_000)
	for range cap(personSlice) {
		personSlice = append(personSlice, person)
	}
	return personSlice
}

func main() {
	myPerson := Person{
		firstName: "Bob",
		lastName:  "Manybob",
		age:       42,
	}
	// mySlice := MakeManyPersons(myPerson)
	mySlice := MakeManyPersonsCap(myPerson)
	fmt.Println(mySlice[len(mySlice)-1])
}
