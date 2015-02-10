package cfsubmit

import (
	"encoding/json"
	"testing"
)

func TestProblemFromJSON(t *testing.T) {
	jsonBytes := []byte(`{"contestId":489,"index":"B","name":"BerSU Ball","type":"PROGRAMMING","points":1000.0,"tags":["dfs and similar","dp","graph matchings","greedy","sortings","two pointers"]}`)
	var p Problem
	if err := json.Unmarshal(jsonBytes, &p); err != nil {
		t.FailNow()
	}
	if p.ContestID != 489 || p.Index != "B" || p.Name != "BerSU Ball" || p.Type != ProblemTypeProgramming || p.Points != 1000.0 {
		t.FailNow()
	}
	tags := []string{"dfs and similar", "dp", "graph matchings", "greedy", "sortings", "two pointers"}
	for i := range tags {
		if tags[i] != p.Tags[i] {
			t.FailNow()
		}
	}
}

func TestUserStatus(t *testing.T) {
	s, err := UserStatus("Fefer_Ivan", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) == 0 || s[0].ID != 9731254 {
		t.Fatal(s[0].ID)
	}
}
