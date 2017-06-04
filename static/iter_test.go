package static

import (
	"context"
	"testing"
)

func TestIterIDs(t *testing.T) {
	iditr := NewIterID(NAMountains)

	ctx := context.Background()

	itrchan := iditr.NextChan(ctx)
	for x := 0; x < 10; x++ {
		s := <-itrchan
		t.Logf("%s", s)
		if x == 0 {
			if s != "Denali" {
				t.Errorf("'Denali' is the first mountain that should be returned")
			}
		}
		if x == 3 {
			if s != "MtStElias" {
				t.Errorf("'MtStElias' should be the fourth name retruned")
			}
		}
	}
}
func TestIterIDs2(t *testing.T) {
	iditr := NewIterID(NAMountains)

	ctx, cancl := context.WithCancel(context.Background())

	itrchan := iditr.NextChan(ctx)
	for x := 0; x < 10; x++ {
		if x > 4 {
			cancl()
		}
		s := <-itrchan
		t.Logf("[%d]: %v", x, s)
		if x == 0 {
			if s != "Denali" {
				t.Errorf("'Denali' is the first mountain that should be returned")
			}
		}
		if x == 3 {
			if s != "MtStElias" {
				t.Errorf("'MtStElias' should be the fourth name retruned")
			}
		}
		if x == 7 {
			if s != "" {
				t.Errorf("Iterator was closed, empty strings should only be returned")
			}
		}
	}
}

func TestIterIDs3(t *testing.T) {
	iditr := NewIterID(NAMountains)

	ctx, cancl := context.WithCancel(context.Background())

	itrchan := iditr.NextChan(ctx)
	for x := 0; x < 110; x++ {
		s := <-itrchan
		if x > 95 {
			t.Logf("[%d]: %v", x, s)
		}
		if x == 100 {
			if s != "Denali" {
				t.Errorf("Denali should be returned after a full name cycle.")
			}
		}
	}
	cancl()
}
