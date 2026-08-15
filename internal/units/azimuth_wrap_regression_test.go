package units

import "testing"

// TestNormalizeAzimuthWrapsBearingsJustWestOfNorth pins the documented output
// domain of NormalizeAzimuth for bearings that land slightly below zero, which
// is what a compass reading close to north turns into once a negative
// declination or compass correction has been applied.
func TestNormalizeAzimuthWrapsBearingsJustWestOfNorth(t *testing.T) {
	cases := []struct {
		input float64
		want  float64
	}{
		{-0.1, 359.9},
		{-0.4, 359.6},
		{-0.75, 359.25},
		{-1, 359},
		{-1.5, 358.5},
		{-2.5, 357.5},
		{-45, 315},
		{-370, 350},
	}
	for _, sample := range cases {
		got := NormalizeAzimuth(sample.input)
		if !NearlyEqual(got, sample.want, 1e-9) {
			t.Fatalf("NormalizeAzimuth(%v) = %v, want %v", sample.input, got, sample.want)
		}
	}
	// Invariants the rest of the pipeline relies on: the result stays inside
	// [0, 360), it names the very same bearing as the input, and it keeps the
	// reciprocal relation with OppositeAzimuth.
	for _, input := range []float64{-0.05, -0.2, -0.9, -3.25, -12.5, -180.5} {
		got := NormalizeAzimuth(input)
		if got < 0 || got >= FullCircleDeg {
			t.Fatalf("NormalizeAzimuth(%v) = %v, which is outside [0, 360)", input, got)
		}
		if !NearlyEqual(got, input+FullCircleDeg, 1e-9) {
			t.Fatalf("NormalizeAzimuth(%v) = %v, want %v", input, got, input+FullCircleDeg)
		}
		if separation := AzimuthSeparation(input, got); separation > 1e-9 {
			t.Fatalf("NormalizeAzimuth(%v) = %v moved the bearing by %v deg", input, got, separation)
		}
		if !NearlyEqual(OppositeAzimuth(input), OppositeAzimuth(got), 1e-9) {
			t.Fatalf("OppositeAzimuth disagrees for %v: %v vs %v", input, OppositeAzimuth(input), OppositeAzimuth(got))
		}
	}
}
