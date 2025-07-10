package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Result struct {
	Id    string
	Value int
}

type Input struct {
	Id   string
	Op   string
	Val1 int
	Val2 int
}

func parser(data []byte) (Input, error) {
	// parse the data
	lines := bytes.Split(data, []byte("\n"))
	// each entry is line 1 id, line 2 operator, line 3 num 1, line 4 num2
	if len(lines) != 4 {
		return Input{}, errors.New("wrong number of lines")
	}
	for _, line := range lines {
		if len(line) == 0 || len(bytes.TrimSpace(line)) == 0 {
			return Input{}, errors.New("line empty")
		}
		if len(line) > 1 && line[0] == '0' {
			return Input{}, errors.New("leading zeros not allowed, except single digit")
		}
		for _, b := range line {
			if b > 127 {
				return Input{}, errors.New("non-ASCII input not allowed")
			}
		}

	}
	id := string(lines[0])
	if len(id) > 61 {
		return Input{}, errors.New("too long")
	}
	op := string(lines[1])
	if len(op) > 1 {
		return Input{}, errors.New("too long")
	}
	if len(string(lines[2])) > 19 {
		return Input{}, errors.New("too long")
	}
	val1, err := strconv.Atoi(string(lines[2]))
	if err != nil {
		return Input{}, err
	}
	if len(string(lines[3])) > 19 {
		return Input{}, errors.New("too long")
	}
	val2, err := strconv.Atoi(string(lines[3]))
	if err != nil {
		return Input{}, err
	}
	if val2 == 0 && op == "/" {
		return Input{}, errors.New("division by zero")
	}
	return Input{
		Id:   id,
		Op:   op,
		Val1: val1,
		Val2: val2,
	}, nil
}

func DataProcessor(in <-chan []byte, out chan<- Result) {
	for data := range in {
		input, err := parser(data)
		if err != nil {
			out <- Result{}
			continue
		}
		var calc int
		switch input.Op {
		case "+":
			calc = input.Val1 + input.Val2
		case "-":
			calc = input.Val1 - input.Val2
		case "*":
			calc = input.Val1 * input.Val2
		case "/":
			calc = input.Val1 / input.Val2
		default:
			out <- Result{}
			continue
		}
		// sum numbers in the data
		// write to another channel
		result := Result{
			Id:    input.Id,
			Value: calc,
		}
		out <- result
	}
	close(out)
}

func WriteData(out <-chan Result, w io.Writer) {
	for r := range out {
		// write the output data to writer
		// each line is id:result
		if r == (Result{}) {
			continue
		}
		w.Write([]byte(fmt.Sprintf("%s:%d\n", r.Id, r.Value)))
	}
}

func NewController(in chan []byte) http.Handler {
	var mu sync.Mutex
	var numSent int
	var numRejected int
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		numSent++
		mu.Unlock()
		// take in data
		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Input"))
			return
		}
		// write it to the queue in raw format
		select {
		case in <- data:
			// success!
		default:
			// if the channel is backed up, return an error
			mu.Lock()
			numRejected++
			mu.Unlock()
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Too Busy: " + strconv.Itoa(numRejected)))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK: " + strconv.Itoa(numSent)))
	})
}

func main() {
	// set everything up
	in := make(chan []byte, 100)
	out := make(chan Result, 100)
	go DataProcessor(in, out)
	f, err := os.Create("results.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	go WriteData(out, f)
	err = http.ListenAndServe(":8080", NewController(in))
	if err != nil {
		fmt.Println(err)
	}
}
