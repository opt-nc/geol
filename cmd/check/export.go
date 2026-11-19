package check

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
)

// exportCheck converts the JSON check output with an in-memory DuckDB session
// and writes the result to stdout.
func exportCheck(jsonData []byte, ext, copyOptions string) error {
	dir, err := os.MkdirTemp("", "geol-check-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	in := filepath.Join(dir, "check.json")
	out := filepath.Join(dir, "check."+ext)
	if err := os.WriteFile(in, jsonData, 0o600); err != nil {
		return err
	}

	db, err := sql.Open("duckdb", "")
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	// A single connection keeps the in-memory database and its view alive.
	db.SetMaxOpenConns(1)

	quote := func(p string) string { return strings.ReplaceAll(p, "'", "''") }
	view := fmt.Sprintf(`CREATE OR REPLACE VIEW V_GEOL_CHECK AS
  SELECT
    t.title,
    s.unnest.software,
    s.unnest.version,
    TRY_CAST(s.unnest.eol_date AS DATE) AS eol_date,
    s.unnest.status,
    TRY_CAST(s.unnest.days AS INTEGER) AS days,
    s.unnest.is_latest,
    s.unnest.latest_version
  FROM read_json('%s') AS t, UNNEST(t.software_components) AS s`, quote(in))
	if _, err := db.Exec(view); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf("COPY V_GEOL_CHECK TO '%s' %s", quote(out), copyOptions)); err != nil {
		return err
	}

	data, err := os.ReadFile(out)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}
