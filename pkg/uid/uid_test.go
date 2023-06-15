package uid

import "testing"

func TestNewID(t *testing.T) {
	t.Run("can generate ID of a certain length", func(t *testing.T) {
		wantLength := 10
		id, err := NewID(10)
		if err != nil {
			t.Fatal(err)
		}
		gotLength := len(id)
		if gotLength != wantLength {
			t.Fatalf("want length %d but got %d", wantLength, gotLength)
		}
	})
}
