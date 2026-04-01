package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Levenshtein Distance ─────────────────────────────────

func TestLevenshtein_IdenticalStrings(t *testing.T) {
	assert.Equal(t, 0, Levenshtein("SSH12345", "SSH12345"))
}

func TestLevenshtein_SingleInsert(t *testing.T) {
	assert.Equal(t, 1, Levenshtein("SSH1234", "SSH12345"), "one insertion → distance should be 1")
}

func TestLevenshtein_SingleSubstitution(t *testing.T) {
	// OCR misread: 'S' → '5'
	assert.Equal(t, 1, Levenshtein("SSH12345", "SSH123S5"), "one substitution → distance should be 1")
}

func TestLevenshtein_CompletelyDifferent(t *testing.T) {
	assert.Equal(t, 3, Levenshtein("ABC", "XYZ"), "completely different 3-char strings → distance should be 3")
}

func TestLevenshtein_EmptyStrings(t *testing.T) {
	assert.Equal(t, 5, Levenshtein("", "hello"))
	assert.Equal(t, 5, Levenshtein("hello", ""))
	assert.Equal(t, 0, Levenshtein("", ""))
}

// ─── Similarity ───────────────────────────────────────────

func TestSimilarity_Identical(t *testing.T) {
	assert.Equal(t, 100.0, Similarity("MED001", "MED001"))
}

func TestSimilarity_OneEditAway(t *testing.T) {
	sim := Similarity("MED001", "MED002")
	// len=6, distance=1, similarity = (1 - 1/6)*100 ≈ 83.33
	assert.InDelta(t, 83.33, sim, 2.0, "one edit → similarity should be ~83.33")
}

func TestSimilarity_BothEmpty(t *testing.T) {
	assert.Equal(t, 100.0, Similarity("", ""))
}

// ─── NormalizeCode ────────────────────────────────────────

func TestNormalizeCode_Standard(t *testing.T) {
	prefix, num, err := NormalizeCode("A0021")
	require.NoError(t, err)
	assert.Equal(t, "A", prefix)
	assert.Equal(t, 21, num)
}

func TestNormalizeCode_LowercasePrefix(t *testing.T) {
	prefix, num, err := NormalizeCode("ssh12345")
	require.NoError(t, err)
	assert.Equal(t, "SSH", prefix)
	assert.Equal(t, 12345, num)
}

func TestNormalizeCode_NoPrefix(t *testing.T) {
	prefix, num, err := NormalizeCode("12345")
	require.NoError(t, err)
	assert.Equal(t, "", prefix)
	assert.Equal(t, 12345, num)
}

func TestNormalizeCode_Invalid(t *testing.T) {
	_, _, err := NormalizeCode("ABCDEF")
	assert.Error(t, err, "expected error for code with no digits")
}

// ─── ExactMatch ───────────────────────────────────────────

func TestExactMatch_LeadingZeros(t *testing.T) {
	// A0021 vs A21 should match after normalization
	assert.True(t, ExactMatch("A0021", "A21"), "A0021 and A21 should be exact match (leading zeros)")
}

func TestExactMatch_CaseInsensitive(t *testing.T) {
	assert.True(t, ExactMatch("ssh123", "SSH123"), "ssh123 and SSH123 should match (case-insensitive)")
}

func TestExactMatch_Different(t *testing.T) {
	assert.False(t, ExactMatch("A001", "B001"), "A001 and B001 should NOT match (different prefix)")
}

// ─── SearchInDatabase ─────────────────────────────────────

func TestSearchInDatabase_ExactMatchFirst(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED002", "MED003"}
	results := c.SearchInDatabase("MED002", dbCodes, 50)

	require.NotEmpty(t, results, "expected at least 1 result")
	assert.Equal(t, "MED002", results[0].Code)
	assert.True(t, results[0].IsExact)
}

func TestSearchInDatabase_ThresholdFiltering(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"ABCDEF"}
	// "XYZXYZ" vs "ABCDEF" → distance=6, similarity=0 → below any threshold
	results := c.SearchInDatabase("XYZXYZ", dbCodes, 50)
	assert.Empty(t, results, "expected 0 results below threshold")
}

func TestSearchInDatabase_EmptyDB(t *testing.T) {
	c := NewOCRClient("")
	results := c.SearchInDatabase("MED001", []string{}, 50)
	assert.Empty(t, results, "empty db should return 0 results")
}

func TestSearchInDatabase_SortedDescending(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED009", "MED002"}
	results := c.SearchInDatabase("MED001", dbCodes, 0)

	for i := 1; i < len(results); i++ {
		assert.GreaterOrEqual(t, results[i-1].Similarity, results[i].Similarity,
			"results not sorted descending at index %d", i)
	}
}

// ─── FindBestMatch ────────────────────────────────────────

func TestFindBestMatch_Found(t *testing.T) {
	c := NewOCRClient("")
	dbCodes := []string{"MED001", "MED002", "MED003"}
	best := c.FindBestMatch("MED001", dbCodes, 50)

	require.NotNil(t, best, "expected a best match")
	assert.Equal(t, "MED001", best.Code)
}

func TestFindBestMatch_NoMatch(t *testing.T) {
	c := NewOCRClient("")
	best := c.FindBestMatch("XYZXYZ", []string{"ABCDEF"}, 90)
	assert.Nil(t, best, "expected nil for no match")
}
