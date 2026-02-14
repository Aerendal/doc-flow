package index

import "testing"

func TestDocDiff(t *testing.T) {
	old := New()
	old.Add(&DocumentRecord{DocID: "a", Path: "a.md", Checksum: "111", MetaChecksum: "m1", BodyChecksum: "b1"})
	old.Add(&DocumentRecord{DocID: "b", Path: "b.md", Checksum: "222", MetaChecksum: "m2", BodyChecksum: "b2"})

	newIdx := New()
	newIdx.Add(&DocumentRecord{DocID: "a", Path: "a.md", Checksum: "999", MetaChecksum: "m1", BodyChecksum: "bX"}) // content changed
	newIdx.Add(&DocumentRecord{DocID: "c", Path: "c.md", Checksum: "333", MetaChecksum: "m3", BodyChecksum: "b3"}) // added

	changes := DocDiff(old, newIdx)
	var added, contentChanged, metaChanged, removed int
	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			added++
		case ChangeContent:
			contentChanged++
		case ChangeMetadata:
			metaChanged++
		case ChangeRemoved:
			removed++
		}
	}
	if added != 1 || contentChanged != 1 || removed != 1 || metaChanged != 0 {
		t.Fatalf("got added=%d content=%d meta=%d removed=%d, want 1/1/0/1", added, contentChanged, metaChanged, removed)
	}
}

func TestDocDiffMetadataOnly(t *testing.T) {
	old := New()
	old.Add(&DocumentRecord{DocID: "a", Path: "a.md", MetaChecksum: "m1", BodyChecksum: "b1"})
	newIdx := New()
	newIdx.Add(&DocumentRecord{DocID: "a", Path: "a.md", MetaChecksum: "m2", BodyChecksum: "b1"})

	changes := DocDiff(old, newIdx)
	if len(changes) != 1 || changes[0].Type != ChangeMetadata {
		t.Fatalf("expected metadata change, got %+v", changes)
	}
	if !changes[0].MetaChanged || changes[0].BodyChanged {
		t.Fatalf("flags not set correctly: %+v", changes[0])
	}
}
