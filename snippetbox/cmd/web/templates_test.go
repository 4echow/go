package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// create slice of anonym. structs containing test case name,
	// input to humanDate func and expected output
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	// loop over test cases
	for _, tt := range tests {
		// use t.Run() to run sub-test for each test case, first param is test name
		// second is anonym. func containing the actual test
		t.Run(tt.name, func(t *testing.T) {
			hd := humandDate(tt.tm)

			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
