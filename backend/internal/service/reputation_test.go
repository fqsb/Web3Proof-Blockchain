package service

import "testing"

func TestGrade(t *testing.T) {
	cases := []struct {
		total uint
		want  string
	}{
		{850, "A"},
		{650, "B"},
		{450, "C"},
		{100, "D"},
	}
	for _, c := range cases {
		if got := grade(c.total); got != c.want {
			t.Fatalf("grade(%d) = %s, want %s", c.total, got, c.want)
		}
	}
}

func TestRecalculateProjectScore(t *testing.T) {
	// grade helper tested; full DB test skipped in unit scope
	if grade(500+300+200) != "A" {
		t.Fatal("max score should be A")
	}
}
