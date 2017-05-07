package static

import "testing"

func TestOrderedList(t *testing.T) {
	if NAMountains[0] != "Denali" {
		t.Errorf("'Denali' should be the first Mountain returned")
	}
}
