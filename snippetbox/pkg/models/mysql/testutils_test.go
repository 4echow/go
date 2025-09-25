package mysql

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// establish sql.DB connection pool for test database
	// use `multiStatements=true` param for multiline setup and teardown SQL files
	// in our DSN for multiple SQL statements in one db.Exec() call
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// read setup script and exec
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// return connection pool and anon. func. which reads and executes teardown script
	// and closes the connection pool. Can be assigned and called later after tests complete
	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
