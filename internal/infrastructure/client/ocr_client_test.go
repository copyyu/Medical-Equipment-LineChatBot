package client

import (
	"testing"
)

// ─── Levenshtein Distance ─────────────────────────────────

func TestLevenshtein_IdenticalStrings(t *testing.T) {
	if got := Levenshtein("SSH12345", "SSH12345"); got != 0 {
		t.Errorf("identical strings should have distance 0, got %d", got)
	}
}

func TestLevenshtein_SingleInsert(t *testing.T) {
	if got := Levenshtein("SSH1234", "SSH12345"); got != 1 {
		t.Errorf("one insertion → distance should be 1, got %d", got)
	}
}

func TestLevenshtein_SingleSubstitution(t *testing.T) {
	// OCR misread: 'S' → '5'
	if got := Levenshtein("SSH12345", "SSH123S5"); got != 1 {
		t.Errorf("one substitution → distance should be 1, got %d", got)
	}
}

func TestLevenshtein_CompletelyDifferent(t *testing.T) {
	if got := Levenshtein("ABC", "XYZ"); got != 3 {
		t.Errorf("completely different 3-char strings → distance should be 3, got %d", got)
	}
}

func TestLevenshtein_EmptyStrings(t *testing.T) {
	if got := Levenshtein("", "hello"); got != 5 {
		t.Errorf("empty vs 5-char string → distance should be 5, got %d", got)
	}
	if got := Levenshtein("hello", ""); got != 5 {
		t.Errorf("5-char vs empty string → distance should be 5, got %d", got)
	}
	if got := Levenshtein("", ""); got != 0 {
		t.Errorf("both empty → distance should be 0, got %d", got)
	}
}

// ─── Similarity ───────────────────────────────────────────

func TestSimilarity_Identical(t *testing.T) {
	sim := Similarity("MED001", "MED001")
	if sim != 100.0 {
		t.Errorf("identical strings → similarity should be 100, got %.2f", sim)
	}
}

func TestSimilarity_OneEditAway(t *testing.T) {
	sim := Similarity("MED001", "MED002")
	// len=6, distance=1, similarity = (1 - 1/6)*100 ≈ 83.33
	if sim < 80 || sim > 85 {
		t.Errorf("one edit → similarity should be ~83.33, got %.2f", sim)
	}
}

func TestSimilarity_BothEmpty(t *testing.T) {
	sim := Similarity("", "")
	if sim != 100.0 {
		t.Errorf("both empty → similarity should be 100, got %.2f", sim)
	}
}

// ─── NormalizeCode ────────────────────────────────────────

func TestNormalizeCode_Standard(t *testing.T) {
	prefix, num, err := NormalizeCode("A0021")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prefix != "A" || num != 21 {
		t.Errorf("expected (A, 21), got (%s, %d)", prefix, num)
	}
}

func TestNormalizeCode_LowercasePrefix(t *testing.T) {
	prefix, num, err := NormalizeCode("ssh12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prefix != "SSH" || num != 12345 {
		t.Errorf("expected (SSH, 12345), got (%s, %d)", prefix, num)
	}
}

func TestNormalizeCode_NoPrefix(t *testing.T) {
	prefix, num, err := NormalizeCode("12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prefix != "" || num != 12345 {
		t.Errorf("expected ('', 12345), got (%s, %d)", prefix, num)
	}
}

func TestNormalizeCode_Invalid(t *testing.T) {
	_, _, err := NormalizeCode("ABCDEF")
	if err == nil {
		t.Error("expected error for code with no digits")
	}
}

// ─── ExactMatch ───────────────────────────────────────────

func TestExactMatch_LeadingZeros(t *testing.T) {
	// A0021 vs A21 should match after normalization
	if !ExactMatch("A0021", "A21") {
		t.Error("A0021 and A21 should be exact match (leading zeros)")
	}
}

func TestExactMatch_CaseInsensitive(t *testing.T) {
	if !ExactMatch("ssh123", "SSH123") {
		t.Error("ssh123 and SSH123 should match (case-insensitive)")
	}
}

func TestExactMatch_Different(t *testing.T) {
	if ExactMatch("A001", "B001") {
		t.Error("A001 and B001 should NOT match (different prefix)")
	}
}

// ─── SearchInDatabase ─────────────────────────────────────

func TestSearchInDatabase_ExactMatchFirst(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED002", "MED003"}
	results := c.SearchInDatabase("MED002", dbCodes, 50)

	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
	if results[0].Code != "MED002" || !results[0].IsExact {
		t.Errorf("first result should be exact match MED002, got %+v", results[0])
	}
}

func TestSearchInDatabase_ThresholdFiltering(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"ABCDEF"}
	// "XYZXYZ" vs "ABCDEF" → distance=6, similarity=0 → below any threshold
	results := c.SearchInDatabase("XYZXYZ", dbCodes, 50)
	if len(results) != 0 {
		t.Errorf("expected 0 results below threshold, got %d", len(results))
	}
}

func TestSearchInDatabase_EmptyDB(t *testing.T) {
	c := NewOCRClient("")
	results := c.SearchInDatabase("MED001", []string{}, 50)
	if len(results) != 0 {
		t.Errorf("empty db should return 0 results, got %d", len(results))
	}
}

func TestSearchInDatabase_SortedDescending(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED009", "MED002"}
	results := c.SearchInDatabase("MED001", dbCodes, 0)

	for i := 1; i < len(results); i++ {
		if results[i].Similarity > results[i-1].Similarity {
			t.Errorf("results not sorted descending: [%d]=%.2f > [%d]=%.2f",
				i, results[i].Similarity, i-1, results[i-1].Similarity)
		}
	}
}

// ─── FindBestMatch ────────────────────────────────────────

func TestFindBestMatch_Found(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED002", "MED003"}
	best := c.FindBestMatch("MED001", dbCodes, 50)

	if best == nil {
		t.Fatal("expected a best match")
	}
	if best.Code != "MED001" {
		t.Errorf("expected MED001, got %s", best.Code)
	}
}

func TestFindBestMatch_NoMatch(t *testing.T) {
	c := NewOCRClient("")
	best := c.FindBestMatch("XYZXYZ", []string{"ABCDEF"}, 90)
	if best != nil {
		t.Errorf("expected nil for no match, got %+v", best)
	}
}
