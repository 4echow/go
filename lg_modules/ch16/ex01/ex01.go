package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type minLenStruct struct {
	A int
	B string `minStrlen:"3"`
	C string `minStrlen:"6"`
}

func ValidateStringLength(s any) error {
	sType := reflect.TypeOf(s)
	if sType.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}
	var errorsList []error
	for i := 0; i < sType.NumField(); i++ {
		curField := sType.Field(i)
		if curField.Type.Kind() != reflect.String {
			continue
		}
		curFieldVal, ok := curField.Tag.Lookup("minStrlen")
		if !ok {
			continue
		}
		minLen, err := strconv.Atoi(curFieldVal)
		if err != nil {
			errorsList = append(errorsList, err)
			continue
		}
		sVal := reflect.ValueOf(s)
		curVal := sVal.Field(i)
		curValLen := len(curVal.String())
		if curValLen < minLen {
			lengthErr := fmt.Errorf("field %s value length %d, expected minimum length %d", curField.Name, curValLen, minLen)
			errorsList = append(errorsList, lengthErr)
		}
	}
	if len(errorsList) != 0 {
		return errors.Join(errorsList...)
	}
	return nil
}

func main() {
	testdata := []minLenStruct{
		{A: 1, B: "abcdef", C: "abcdef"},
		{A: 2, B: "ab", C: "abcdef"},
		{A: 3, B: "abc", C: "abc"},
	}
	for _, d := range testdata {
		err := ValidateStringLength(d)
		if err != nil {
			fmt.Println(d, err)
		}
	}
}
