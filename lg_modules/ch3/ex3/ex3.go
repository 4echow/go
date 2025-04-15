package main

import "fmt"

func main() {
	type Employee struct {
		firstName string
		lastName  string
		id        int
	}

	firstEmployee := Employee{
		"Max",
		"Mustermann",
		42,
	}
	secondEmployee := Employee{
		firstName: "Maria",
		lastName:  "Musterfrau",
		id:        24,
	}
	var thirdEmployee struct {
		firstName string
		lastName  string
		id        int
	}
	thirdEmployee.firstName = "Hasso"
	thirdEmployee.lastName = "Musterhund"
	thirdEmployee.id = 5

	fmt.Println(firstEmployee)
	fmt.Println(secondEmployee)
	fmt.Println(thirdEmployee)
}
