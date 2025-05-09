package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidId = errors.New("this id is invalid")

type ErrEmptyField struct {
	employeeField string
}

func (e ErrEmptyField) Error() string {
	return e.employeeField
}

func processError(emp Employee, err error) string {
	var emptyFieldError ErrEmptyField
	if errors.Is(err, ErrInvalidId) {
		return fmt.Sprintf("error: %v, ID: %s", err, emp.ID)
	} else if errors.As(err, &emptyFieldError) {
		return fmt.Sprintf("empty field: %s", emptyFieldError.employeeField)
	} else {
		return fmt.Sprintf("error: %v", err)
	}
}

func main() {
	d := json.NewDecoder(strings.NewReader(data))
	count := 0
	for d.More() {
		count++
		var emp Employee
		err := d.Decode(&emp)
		if err != nil {
			fmt.Printf("record %d: %v\n", count, err)
			continue
		}
		err = ValidateEmployee(emp)
		if err != nil {
			switch e := err.(type) {
			case interface{ Unwrap() []error }:
				var allErrors []string
				for _, innerErr := range e.Unwrap() {
					allErrors = append(allErrors, processError(emp, innerErr))
				}
				fmt.Printf("record %d: validation errors: %s\n", count, strings.Join(allErrors, ", "))
			default:
				processedError := processError(emp, err)
				fmt.Printf("record %d: validation error: %s\n", count, processedError)
			}
		} else {
			fmt.Printf("record %d: valid: %+v\n", count, emp)
		}

	}
}

const data = `
{
	"id": "ABCD-123",
	"first_name": "Bob",
	"last_name": "Bobson",
	"title": "Senior Manager"
}
{
	"id": "XYZ-123",
	"first_name": "Mary",
	"last_name": "Maryson",
	"title": "Vice President"
}
{
	"id": "BOTX-263",
	"first_name": "",
	"last_name": "Garciason",
	"title": "Manager"
}
{
	"id": "HLXO-829",
	"first_name": "Pierre",
	"last_name": "",
	"title": "Intern"
}
{
	"id": "MOXW-821",
	"first_name": "Franklin",
	"last_name": "Watanabe",
	"title": ""
}
{
	"id": "",
	"first_name": "Shelly",
	"last_name": "Shellson",
	"title": "CEO"
}
{
	"id": "YDOD-324",
	"first_name": "",
	"last_name": "",
	"title": ""
}
`

type Employee struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string `json:"title"`
}

var (
	validID = regexp.MustCompile(`\w{4}-\d{3}`)
)

func ValidateEmployee(e Employee) error {
	var errs []error
	if len(e.ID) == 0 {
		errs = append(errs, ErrEmptyField{employeeField: "ID"})
	}
	if !validID.MatchString(e.ID) {
		errs = append(errs, ErrInvalidId)
	}
	if len(e.FirstName) == 0 {
		errs = append(errs, ErrEmptyField{employeeField: "FirstName"})
	}
	if len(e.LastName) == 0 {
		errs = append(errs, ErrEmptyField{employeeField: "LastName"})
	}
	if len(e.Title) == 0 {
		errs = append(errs, ErrEmptyField{employeeField: "Title"})
	}
	if len(errs) > 1 {
		return errors.Join(errs...)
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return nil
}
