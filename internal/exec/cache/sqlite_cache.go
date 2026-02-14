package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const schemaVersion = 1

type sqliteCache struct {
	db *sql.DB
}

func OpenSQLite(path string) (Cache, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("cannot create cache dir: %w", err)
	}

	dsn := path + "?_pragma=journal_mode(wal)&_pragma=foreign_keys(on)&_pragma=synchronous(normal)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite cache: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("cannot ping sqlite cache: %w", err)
	}

	c := &sqliteCache{db: db}
	if err := c.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return c, nil
}

func (c *sqliteCache) Close() error {
	if c == nil || c.db == nil {
		return nil
	}
	return c.db.Close()
}

func (c *sqliteCache) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS cache_meta (
  schema_version INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  tool_version TEXT NOT NULL DEFAULT '',
  note TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS documents (
  path TEXT PRIMARY KEY,
  root_id TEXT NOT NULL DEFAULT '',
  size INTEGER NOT NULL,
  mtime_unix INTEGER NOT NULL,
  checksum_meta TEXT NOT NULL DEFAULT '',
  checksum_body TEXT NOT NULL DEFAULT '',
  checksum_full TEXT NOT NULL DEFAULT '',
  doc_id TEXT NOT NULL DEFAULT '',
  doc_type TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  owner TEXT NOT NULL DEFAULT '',
  language TEXT NOT NULL DEFAULT '',
  priority TEXT NOT NULL DEFAULT '',
  context_sources_json TEXT NOT NULL DEFAULT '[]',
  depends_on_json TEXT NOT NULL DEFAULT '[]',
  heading_count INTEGER NOT NULL DEFAULT 0,
  max_depth INTEGER NOT NULL DEFAULT 0,
  lines INTEGER NOT NULL DEFAULT 0,
  parse_error_code TEXT NOT NULL DEFAULT '',
  parse_error_message TEXT NOT NULL DEFAULT '',
  facts_json TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_documents_doc_id ON documents(doc_id);
CREATE INDEX IF NOT EXISTS idx_documents_checksum_meta ON documents(checksum_meta);
CREATE INDEX IF NOT EXISTS idx_documents_checksum_body ON documents(checksum_body);
CREATE INDEX IF NOT EXISTS idx_documents_mtime ON documents(mtime_unix);
CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);

CREATE TABLE IF NOT EXISTS runs (
  run_id TEXT PRIMARY KEY,
  started_at TEXT NOT NULL,
  finished_at TEXT NOT NULL DEFAULT '',
  mode TEXT NOT NULL DEFAULT '',
  root TEXT NOT NULL DEFAULT '',
  changed_files INTEGER NOT NULL DEFAULT 0,
  cache_hits INTEGER NOT NULL DEFAULT 0,
  cache_misses INTEGER NOT NULL DEFAULT 0
);`
	if _, err := c.db.Exec(schema); err != nil {
		return fmt.Errorf("cache migrate schema failed: %w", err)
	}

	var current int
	err := c.db.QueryRow(`SELECT schema_version FROM cache_meta LIMIT 1`).Scan(&current)
	switch {
	case err == sql.ErrNoRows:
		_, err = c.db.Exec(`INSERT INTO cache_meta(schema_version, created_at) VALUES(?, ?)`, schemaVersion, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("cache init meta failed: %w", err)
		}
		return nil
	case err != nil:
		return fmt.Errorf("cache read meta failed: %w", err)
	}

	if current == schemaVersion {
		return nil
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM documents`); err != nil {
		return fmt.Errorf("cache reset documents failed: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM runs`); err != nil {
		return fmt.Errorf("cache reset runs failed: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM cache_meta`); err != nil {
		return fmt.Errorf("cache reset meta failed: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO cache_meta(schema_version, created_at) VALUES(?, ?)`, schemaVersion, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("cache rewrite meta failed: %w", err)
	}
	return tx.Commit()
}

func (c *sqliteCache) Get(path string, fp Fingerprint) (*DocumentFacts, bool, error) {
	var (
		size         int64
		mtimeUnix    int64
		checksumMeta string
		checksumBody string
		factsJSON    string
	)
	err := c.db.QueryRow(
		`SELECT size, mtime_unix, checksum_meta, checksum_body, facts_json FROM documents WHERE path = ?`,
		path,
	).Scan(&size, &mtimeUnix, &checksumMeta, &checksumBody, &factsJSON)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	if size != fp.Size || mtimeUnix != fp.MTimeUnix {
		return nil, false, nil
	}
	if fp.MetaChecksum != "" && checksumMeta != "" && checksumMeta != fp.MetaChecksum {
		return nil, false, nil
	}
	if fp.BodyChecksum != "" && checksumBody != "" && checksumBody != fp.BodyChecksum {
		return nil, false, nil
	}

	var facts DocumentFacts
	if err := json.Unmarshal([]byte(factsJSON), &facts); err != nil {
		return nil, false, err
	}
	return &facts, true, nil
}

func (c *sqliteCache) Put(path string, fp Fingerprint, facts DocumentFacts) error {
	factsJSONBytes, err := json.Marshal(facts)
	if err != nil {
		return err
	}
	contextJSON, _ := json.Marshal(facts.ContextSources)
	depsJSON, _ := json.Marshal(facts.DependsOn)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = c.db.Exec(
		`INSERT INTO documents(
			path, root_id, size, mtime_unix, checksum_meta, checksum_body, checksum_full,
			doc_id, doc_type, status, version, owner, language, priority,
			context_sources_json, depends_on_json,
			heading_count, max_depth, lines,
			parse_error_code, parse_error_message,
			facts_json, updated_at
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			root_id=excluded.root_id,
			size=excluded.size,
			mtime_unix=excluded.mtime_unix,
			checksum_meta=excluded.checksum_meta,
			checksum_body=excluded.checksum_body,
			checksum_full=excluded.checksum_full,
			doc_id=excluded.doc_id,
			doc_type=excluded.doc_type,
			status=excluded.status,
			version=excluded.version,
			owner=excluded.owner,
			language=excluded.language,
			priority=excluded.priority,
			context_sources_json=excluded.context_sources_json,
			depends_on_json=excluded.depends_on_json,
			heading_count=excluded.heading_count,
			max_depth=excluded.max_depth,
			lines=excluded.lines,
			parse_error_code=excluded.parse_error_code,
			parse_error_message=excluded.parse_error_message,
			facts_json=excluded.facts_json,
			updated_at=excluded.updated_at`,
		path, facts.Root, fp.Size, fp.MTimeUnix, fp.MetaChecksum, fp.BodyChecksum, facts.ChecksumFull,
		facts.DocID, facts.DocType, facts.Status, facts.Version, facts.Owner, facts.Language, facts.Priority,
		string(contextJSON), string(depsJSON),
		facts.HeadingCount, facts.MaxDepth, facts.Lines,
		facts.ParseErrorCode, facts.ParseErrorMessage,
		string(factsJSONBytes), now,
	)
	return err
}

func (c *sqliteCache) BeginRun(mode string, root string) (string, error) {
	runID := fmt.Sprintf("run_%d", time.Now().UTC().UnixNano())
	_, err := c.db.Exec(
		`INSERT INTO runs(run_id, started_at, mode, root) VALUES (?, ?, ?, ?)`,
		runID, time.Now().UTC().Format(time.RFC3339), mode, root,
	)
	if err != nil {
		return "", err
	}
	return runID, nil
}

func (c *sqliteCache) EndRun(runID string, stats RunStats) error {
	_, err := c.db.Exec(
		`UPDATE runs SET finished_at = ?, changed_files = ?, cache_hits = ?, cache_misses = ? WHERE run_id = ?`,
		time.Now().UTC().Format(time.RFC3339), stats.ChangedFiles, stats.CacheHits, stats.CacheMisses, runID,
	)
	return err
}
