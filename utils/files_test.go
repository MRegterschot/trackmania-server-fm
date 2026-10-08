package utils

import "testing"

func TestDecodeUploadFilename(t *testing.T) {
	cases := map[string]string{
		"5 - Go %22Snow%22 Car.Map.Gbx": `5 - Go "Snow" Car.Map.Gbx`,
		"plain.Map.Gbx":                 "plain.Map.Gbx",
		"line%0D%0Abreak.Map.Gbx":       "linebreak.Map.Gbx",
		" 3 - PipeLand.Map.Gbx":         " 3 - PipeLand.Map.Gbx",
	}
	for in, want := range cases {
		if got := DecodeUploadFilename(in); got != want {
			t.Errorf("DecodeUploadFilename(%q) = %q, want %q", in, got, want)
		}
	}
}
