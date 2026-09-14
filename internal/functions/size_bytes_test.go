package functions

import "testing"

func TestParseSize(t *testing.T) {
	for in, want := range map[string]int64{
		"0": 0, "512": 512, "1K": 1024, "1k": 1024, "10G": 10 << 30, "1.5 TiB": 3 << 39, "2GB": 2 << 30,
		" 500M ": 500 << 20, "1P": 1 << 50, "0.5K": 512,
	} {
		got, err := ParseSize(in)
		if err != nil || got != want {
			t.Errorf("ParseSize(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	for _, in := range []string{"", "ten", "1X", "-1G", "1.1", "9999999P"} {
		if _, err := ParseSize(in); err == nil {
			t.Errorf("ParseSize(%q) accepted", in)
		}
	}
}
