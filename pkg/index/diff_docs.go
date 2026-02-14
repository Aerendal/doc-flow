package index

type DocChangeType string

const (
	ChangeAdded        DocChangeType = "added"
	ChangeRemoved      DocChangeType = "removed"
	ChangeContent      DocChangeType = "content_changed"
	ChangeMetadata     DocChangeType = "metadata_changed"
	ChangeChecksumOnly DocChangeType = "changed"
)

type DocChange struct {
	DocID       string
	Path        string
	Type        DocChangeType
	MetaChanged bool
	BodyChanged bool
}

// DocDiff compares old vs new index and returns changed/added/removed docs (non-templates included).
func DocDiff(oldIdx, newIdx *DocumentIndex) []DocChange {
	oldMap := make(map[string]*DocumentRecord)
	for _, r := range oldIdx.All() {
		if r.DocID != "" {
			oldMap[r.DocID] = r
		}
	}
	newMap := make(map[string]*DocumentRecord)
	for _, r := range newIdx.All() {
		if r.DocID != "" {
			newMap[r.DocID] = r
		}
	}

	var changes []DocChange

	// added/changed
	for id, n := range newMap {
		o, ok := oldMap[id]
		if !ok {
			changes = append(changes, DocChange{DocID: id, Path: n.Path, Type: ChangeAdded})
			continue
		}
		metaChanged := o.MetaChecksum != n.MetaChecksum
		bodyChanged := o.BodyChecksum != n.BodyChecksum
		switch {
		case bodyChanged:
			changes = append(changes, DocChange{DocID: id, Path: n.Path, Type: ChangeContent, MetaChanged: metaChanged, BodyChanged: true})
		case metaChanged:
			changes = append(changes, DocChange{DocID: id, Path: n.Path, Type: ChangeMetadata, MetaChanged: true})
		case o.Checksum != n.Checksum:
			changes = append(changes, DocChange{DocID: id, Path: n.Path, Type: ChangeChecksumOnly})
		}
	}
	// removed
	for id, o := range oldMap {
		if _, ok := newMap[id]; !ok {
			changes = append(changes, DocChange{DocID: id, Path: o.Path, Type: ChangeRemoved})
		}
	}

	return changes
}
