package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// init new httptest.ResponseRecorder and dummy http.Request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a mock HTTP handler we can pass to our secureHeaders
	// middleware, which writes 200 status code and "OK" response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock handler to secureHeaders middleware, since it *returns*
	// a http.Handler, we can call its ServeHTTP() method, passing the
	// http.ResponseRecorder and dummy http.Request to execute it
	secureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on http.ResponseRecorder to get the results
	rs := rr.Result()

	// check that middleware correctly set X-Frame-Options header on response
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// check for correct X-XSS-Protection header
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	// check that middleware correctly called next handler in line and
	// response status and body are as expected
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
