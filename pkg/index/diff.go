package index

// ChangedTemplates returns doc_ids of templates whose checksum changed between two indexes.
func ChangedTemplates(oldIdx, newIdx *DocumentIndex) map[string]bool {
	changed := make(map[string]bool)
	oldMap := make(map[string]*DocumentRecord)
	for _, r := range oldIdx.All() {
		if r.DocID != "" {
			oldMap[r.DocID] = r
		}
	}
	for _, r := range newIdx.All() {
		if r.DocID == "" {
			continue
		}
		if r.DocType != "template" && r.TemplateSource == "" {
			continue
		}
		if old, ok := oldMap[r.DocID]; ok {
			if old.Checksum != r.Checksum {
				changed[r.DocID] = true
			}
		} else {
			changed[r.DocID] = true
		}
	}
	return changed
}

// Impact computes docs affected by changed templates via TemplateSource match.
func Impact(changedTemplates map[string]bool, idx *DocumentIndex) map[string][]string {
	result := make(map[string][]string)
	for _, r := range idx.All() {
		if r.TemplateSource == "" {
			continue
		}
		if changedTemplates[r.TemplateSource] {
			result[r.TemplateSource] = append(result[r.TemplateSource], r.Path)
		}
	}
	return result
}
