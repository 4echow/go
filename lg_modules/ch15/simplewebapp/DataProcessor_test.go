package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) {
	return 0, errors.New("simulated read failure")
}
func (brokenReader) Close() error {
	return nil
}

func TestParser(t *testing.T) {
	data := []struct {
		name       string
		expression string
		expected   Input
		errMsg     string
	}{
		{
			"Addition",
			"Addition\n+\n2\n3",
			Input{
				Id:   "Addition",
				Op:   "+",
				Val1: 2,
				Val2: 3,
			},
			"",
		},
		{
			"Subtraction",
			"Subtraction\n-\n4\n7",
			Input{
				Id:   "Subtraction",
				Op:   "-",
				Val1: 4,
				Val2: 7,
			},
			"",
		},
		{
			"Multiplication",
			"Multiplication\n*\n3\n7",
			Input{
				Id:   "Multiplication",
				Op:   "*",
				Val1: 3,
				Val2: 7,
			},
			"",
		},
		{
			"Division_valid",
			"Division_valid\n/\n6\n3",
			Input{
				Id:   "Division_valid",
				Op:   "/",
				Val1: 6,
				Val2: 3,
			},
			"",
		},
		{
			"Division_invalid",
			"Division_invalid\n/\n3\n0",
			Input{},
			"division by zero",
		},
		{
			"Invalid_Val1",
			"Invalid_Val1\n/\nfoo\n7",
			Input{},
			`strconv.Atoi: parsing "foo": invalid syntax`,
		},
		{
			"Invalid_Val2",
			"Invalid_Val2\n/\n7\nbar",
			Input{},
			`strconv.Atoi: parsing "bar": invalid syntax`,
		},
		{
			"Invalid_Linenum",
			"Invalid_Linenum1\n/\n7\n6\nfoo",
			Input{},
			"wrong number of lines",
		},
		{
			"Empty",
			"\n\n\n",
			Input{},
			"line empty",
		},
		{
			"Invalid_Characters",
			"0\n+\n2\n\xff",
			Input{},
			"non-ASCII input not allowed",
		},
		{
			"Too_long_Id",
			"12345678901234567890123456789012345678901234567890123456789012\n%\n3\n5",
			Input{},
			"too long",
		},
		{
			"Too_long_Op",
			"Too_long_Op\n++\n1\n2",
			Input{},
			"too long",
		},
		{
			"Too_long_Var1",
			"Too_long_Var1\n+\n12345678901234567890\n1",
			Input{},
			"too long",
		},
		{
			"Too_long_Var2",
			"Too_long_Var2\n+\n7\n12345678901234567890",
			Input{},
			"too long",
		},
		{
			"Leading_zeros_not_allowed",
			"Leading_zeros_not_allowed\n-\n01235\n7",
			Input{},
			"leading zeros not allowed, except single digit",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			input, err := parser([]byte(d.expression))
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if diff := cmp.Diff(d.expected, input); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(d.errMsg, errMsg); diff != "" {
				t.Error(diff)
			}
		})
	}

}

func ToData(in Input) []byte {
	return []byte(fmt.Sprintf("%s\n%s\n%d\n%d", in.Id, in.Op, in.Val1, in.Val2))
}

func FuzzParser(f *testing.F) {
	testcases := [][]byte{
		[]byte("Addition\n+\n2\n3"),
		[]byte("Subtraction\n-\n7\n4"),
		[]byte("Multiplication\n*\n7\n3"),
		[]byte("Division\n/\n9\n3"),
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		input, err := parser(a)
		if err != nil {
			t.Skip("handled error")
		}
		roundTrip := ToData(input)
		input2, err := parser(roundTrip)
		if err != nil {
			t.Errorf("roundtrip parser error: %v\ndata: %q", err, roundTrip)
		}
		if diff := cmp.Diff(input, input2); diff != "" {
			t.Error(diff)
		}
	})
}

func TestDataProcessor(t *testing.T) {
	data := []struct {
		name       string
		expression string
		expected   Result
	}{
		{
			"Addition",
			"Addition\n+\n2\n3",
			Result{
				Id:    "Addition",
				Value: 5,
			},
		},
		{
			"Subtraction",
			"Subtraction\n-\n4\n7",
			Result{
				Id:    "Subtraction",
				Value: -3,
			},
		},
		{
			"Multiplication",
			"Multiplication\n*\n3\n7",
			Result{
				Id:    "Multiplication",
				Value: 21,
			},
		},
		{
			"Division_valid",
			"Division_valid\n/\n6\n3",
			Result{
				Id:    "Division_valid",
				Value: 2,
			},
		},
		{
			"Division_invalid",
			"Division_invalid\n/\n3\n0",
			Result{},
		},
		{
			"Unknown",
			"Unknown\n%\n3\n5",
			Result{},
		},
	}
	in := make(chan []byte)
	out := make(chan Result)

	go DataProcessor(in, out)

	go func() {
		for _, d := range data {
			in <- []byte(d.expression)
		}
		close(in)
	}()
	for _, d := range data {
		d := d // shadowing just in case
		t.Run(d.name, func(t *testing.T) {
			if diff := cmp.Diff(d.expected, <-out); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestWriteData(t *testing.T) {
	data := []struct {
		name     string
		input    Result
		expected string
	}{
		{
			name:     "success_write1",
			input:    Result{Id: "abc", Value: 42},
			expected: "abc:42\n",
		},
		{
			name:     "success_write2",
			input:    Result{Id: "def", Value: -3},
			expected: "def:-3\n",
		},
		{
			name:     "failure_write1",
			input:    Result{},
			expected: "",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			ch := make(chan Result, 1)
			var buf bytes.Buffer

			ch <- d.input
			close(ch)

			WriteData(ch, &buf)

			if diff := cmp.Diff(d.expected, buf.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestNewController(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		in := make(chan []byte, 1)
		handler := NewController(in)
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Post(server.URL, "text/plain", strings.NewReader("Addition\n+\n5\n3"))
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		if http.StatusAccepted != resp.StatusCode {
			t.Errorf("Expected Status %d got %d", http.StatusAccepted, resp.StatusCode)
		}

		bodybytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		body := string(bodybytes)

		if diff := cmp.Diff(body, "OK: 1"); diff != "" {
			t.Error(diff)
		}
		channelInput := <-in
		if diff := cmp.Diff("Addition\n+\n5\n3", string(channelInput)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("Service_Unavailable", func(t *testing.T) {
		in := make(chan []byte, 1)
		in <- []byte("create blocked channel")
		handler := NewController(in)
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Post(server.URL, "text/plain", strings.NewReader("Addition\n+\n5\n3"))
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		if http.StatusServiceUnavailable != resp.StatusCode {
			t.Errorf("Expected Status %d got %d", http.StatusServiceUnavailable, resp.StatusCode)
		}

		bodybytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		body := string(bodybytes)

		if diff := cmp.Diff(body, "Too Busy: 1"); diff != "" {
			t.Error(diff)
		}
		channelInput := <-in
		if diff := cmp.Diff("create blocked channel", string(channelInput)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("Bad_Input", func(t *testing.T) {

		in := make(chan []byte, 1)
		handler := NewController(in)

		// Create fake request with broken body
		req := httptest.NewRequest(http.MethodPost, "/", brokenReader{})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "Bad Input") {
			t.Errorf("expected response body to contain 'Bad Input', got %q", string(body))
		}
	})
}
