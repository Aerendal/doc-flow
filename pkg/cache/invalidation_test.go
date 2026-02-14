package cache

import (
	"testing"

	"docflow/pkg/index"
)

func TestInvalidateDependents(t *testing.T) {
	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "a", DependsOn: []string{"b"}})
	idx.Add(&index.DocumentRecord{DocID: "b", DependsOn: []string{"c"}})
	idx.Add(&index.DocumentRecord{DocID: "c"})
	rev := BuildReverseDeps(idx)

	toInv := InvalidateDependents([]string{"c"}, rev)
	want := map[string]bool{"c": true, "b": true, "a": true}
	if len(toInv) != 3 {
		t.Fatalf("got %v", toInv)
	}
	for _, id := range toInv {
		if !want[id] {
			t.Fatalf("unexpected id %s", id)
		}
	}
}
