package static

import "testing"

func TestOrderedList(t *testing.T) {
	if NAMountains[0] != "Denali" {
		t.Errorf("'Denali' should be the first Mountain returned")
	}
}

func TestStringNormalize(t *testing.T) {
	a := "VolcánAcatenango"
	b := "VolcanAcatenango"
	normalized, err := normalizeString(a)
	if err != nil {
		t.Errorf("error converting %q to %q: %v", a, b, err)
	}
	if normalized != b {
		t.Errorf("normalized %q does not match expected %q", normalized, b)
	}
	t.Logf("b: %q", b)
}

func TestNaNormalized(t *testing.T) {
	normed := NormalizedNaMountains()
	if len(normed) != len(NAMountains) {
		t.Errorf("normalized NA Mountain length do not match: %d:%d", len(normed), len(NAMountains))
	}
}
