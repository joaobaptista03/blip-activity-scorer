package main

import (
	"testing"
)

func TestDeduplicator(t *testing.T) {
	deduper := NewDeduplicator()

	c1 := Commit{Timestamp: 100, Username: "user1", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3}
	c2 := Commit{Timestamp: 100, Username: "user1", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3} // Exact duplicate
	c3 := Commit{Timestamp: 101, Username: "user1", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3} // Different timestamp
	c4 := Commit{Timestamp: 100, Username: "user2", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3} // Different user
	c5 := Commit{Timestamp: 100, Username: "", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3}      // Blank username (distinct from "user1")
	c6 := Commit{Timestamp: 100, Username: "", Repository: "repo1", Files: 1, Additions: 2, Deletions: 3}      // Duplicate blank username

	if deduper.IsDuplicate(c1) {
		t.Error("c1 should not be a duplicate on first sighting")
	}
	if !deduper.IsDuplicate(c2) {
		t.Error("c2 should be a duplicate of c1")
	}
	if deduper.IsDuplicate(c3) {
		t.Error("c3 should not be a duplicate (different timestamp)")
	}
	if deduper.IsDuplicate(c4) {
		t.Error("c4 should not be a duplicate (different user)")
	}
	if deduper.IsDuplicate(c5) {
		t.Error("c5 should not be a duplicate (first blank username)")
	}
	if !deduper.IsDuplicate(c6) {
		t.Error("c6 should be a duplicate of c5 (matching blank username)")
	}
}
