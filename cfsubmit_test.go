package cfsubmit

import "testing"

func TestCfsNewCorrect(t *testing.T) {
	filename := "123a.cpp"
	s, err := New(filename)
	if err != nil {
		t.FailNow()
	}
	if s.ContestID != "123" || s.ProblemID != "A" || s.Extension != "cpp" {
		t.FailNow()
	}
}

func TestCfsNewBadcID(t *testing.T) {
	filename := "a123a.cpp"
	_, err := New(filename)
	if err == nil {
		t.FailNow()
	}
}

func TestCfsNewNoExt(t *testing.T) {
	filename := "123a"
	_, err := New(filename)
	if err == nil {
		t.FailNow()
	}
}
