package cache

import (
	"path/filepath"
	"testing"
)

func TestSQLiteCachePutGetHitAndMiss(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cache.sqlite")
	c, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	fp := Fingerprint{Size: 10, MTimeUnix: 12345, MetaChecksum: "m1", BodyChecksum: "b1"}
	in := DocumentFacts{
		Path:      "a.md",
		Root:      "/tmp/root",
		Size:      10,
		MTimeUnix: 12345,
		DocID:     "a",
		Content:   "---\n---\nbody",
	}
	if err := c.Put("a.md", fp, in); err != nil {
		t.Fatal(err)
	}

	got, hit, err := c.Get("a.md", fp)
	if err != nil {
		t.Fatal(err)
	}
	if !hit {
		t.Fatal("expected cache hit")
	}
	if got == nil || got.DocID != "a" {
		t.Fatalf("unexpected facts: %#v", got)
	}

	_, hit, err = c.Get("a.md", Fingerprint{Size: 11, MTimeUnix: 12345})
	if err != nil {
		t.Fatal(err)
	}
	if hit {
		t.Fatal("expected cache miss for different fingerprint")
	}
}

func TestSQLiteCacheSchemaVersionReset(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cache.sqlite")
	c, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Put("a.md", Fingerprint{Size: 1, MTimeUnix: 1}, DocumentFacts{Path: "a.md", DocID: "a"}); err != nil {
		t.Fatal(err)
	}

	impl, ok := c.(*sqliteCache)
	if !ok {
		t.Fatal("unexpected cache implementation")
	}
	if _, err := impl.db.Exec(`UPDATE cache_meta SET schema_version = 999`); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()

	c2, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()

	_, hit, err := c2.Get("a.md", Fingerprint{Size: 1, MTimeUnix: 1})
	if err != nil {
		t.Fatal(err)
	}
	if hit {
		t.Fatal("expected cache reset after schema mismatch")
	}
}
