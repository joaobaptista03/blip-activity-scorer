package main

// Deduplicator keeps track of seen commits to filter exact duplicate rows.
type Deduplicator struct {
	seen map[Commit]struct{}
}

// NewDeduplicator creates and returns a new Deduplicator.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		seen: make(map[Commit]struct{}),
	}
}

// IsDuplicate checks if the commit has been processed before. If not, it registers it and returns false.
func (d *Deduplicator) IsDuplicate(c Commit) bool {
	if _, exists := d.seen[c]; exists {
		return true
	}
	d.seen[c] = struct{}{}
	return false
}
