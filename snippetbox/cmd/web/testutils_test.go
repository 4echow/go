package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/4echow/go/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

// define regexp that captures csrf token value from HTML for user signup page
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

// extractCSRFToken uses FindSubmatch method to extract token from html body
// returns array with entire matched pattern in first position, values of captured data
// in subsequent positions
func extractCSRFToken(t *testing.T, body []byte) string {
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	// needs unescaping since html/template potentially escapes characters like "+"
	return html.UnescapeString(string(matches[1]))
}

// newTestApplication returns an instance of our application struct
// containing mocked dependencies
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// create session manager instance, same settings as production
	session := sessions.New([]byte("738d4b12617a007521567fffa4e811673cd35bdac09bd716270e0d7bc9fe7be3"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// init dependencies using the mocks for loggers and database models
	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
	}
}

// testServer anonymously embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// newTestServer initializes and returns a new instance of our
// custom testServer type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// init a cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// add jar to client, so response cookies are stored and sent
	// with subsequent requests
	ts.Client().Jar = jar

	// disable redirect-following for client
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get method on our custom testServer type. This method makes a GET
// request to a given url path on the test server, and returns the
// response status code, headers and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL + urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// read response body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	// return response status, headers and body
	return rs.StatusCode, rs.Header, body
}
