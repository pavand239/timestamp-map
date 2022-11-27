package timestampmap

import "testing"

var testSliceForSearchClosest = []int64{1, 3, 5, 6, 7, 8, 10, 12, 23, 24, 25, 27, 31, 42, 44, 56, 100, 399}

var testResMapForSearchClosest = map[int64]int64{
	-1:   1,
	7:    7,
	9:    8,
	18:   23,
	26:   25,
	33:   31,
	39:   42,
	100:  100,
	250:  399,
	1000: 399,
}

func TestSearchClosest(t *testing.T) {
	for needle, value := range testResMapForSearchClosest {
		t.Logf("Search %v", needle)
		foundValue, err := searchClosestKey(testSliceForSearchClosest, needle)
		if err != nil {
			t.Fatal(err)
		}

		if foundValue != value {
			t.Fatalf("Found %v, want %v", foundValue, value)
		} else {
			t.Logf("Found expected value %v", foundValue)
		}
	}
}
