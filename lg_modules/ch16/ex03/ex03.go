package main

/*
	#include <string.h>

	int mini_calc(char *op, int a, int b) {
		if (strcmp(op, "+") == 0) {
			return a + b;
		}
		if (strcmp(op, "*") == 0) {
			return a * b;
		}
		if (strcmp(op, "-") == 0) {
			return a - b;
		}
		if (strcmp(op, "/") == 0) {
			if (b == 0) {
				return 0;
			}
			return a / b;
		}
		return 0;
	}
*/
import "C"

import (
	"fmt"
)

func main() {
	type teststruct struct {
		Op   string
		Var1 int
		Var2 int
	}
	data := []teststruct{
		{
			Op:   "+",
			Var1: 5,
			Var2: 7,
		},
		{
			Op:   "-",
			Var1: 5,
			Var2: 7,
		},
		{
			Op:   "*",
			Var1: 5,
			Var2: 7,
		},
		{
			Op:   "/",
			Var1: 14,
			Var2: 7,
		},
		{
			Op:   "/",
			Var1: 5,
			Var2: 7,
		},
		{
			Op:   "/",
			Var1: 5,
			Var2: 0,
		},
	}
	for _, d := range data {
		cOp := C.CString(d.Op)
		cVar1 := C.int(d.Var1)
		cVar2 := C.int(d.Var2)
		result := C.mini_calc(cOp, cVar1, cVar2)
		fmt.Printf("Result for '%s' operation with values %d and %d: %d\n", d.Op, d.Var1, d.Var2, result)
	}
}
