package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/4echow/go/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

// newTestApplication returns an instance of our application struct
// containing mocked dependencies
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// create session manager instance, same settings as production
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// init dependencies using the mocks for loggers and database models
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		session: session,
		snippets: &mock.SnippetModel{},
		templateCache: templateCache,
		users: &mock.UserModel{},
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
