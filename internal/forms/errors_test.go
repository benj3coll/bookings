package forms

import "testing"

func TestAdd(t *testing.T) {
	var e errors = make(map[string][]string)
	e.Add("p1", "Test p1")
	if len(e) != 1 {
		t.Error("errors zu gross oder zu klein")
	}
}

func TestGet(t *testing.T) {
	var e errors = make(map[string][]string)
	e.Add("p1", "Test p1")
	if e.Get("p1") != "Test p1" {
		t.Error("p1 hat falschen Wert")
	}
	if e.Get("pnf") != "" {
		t.Error("pnf darf nicht gefunden werden")
	}
}
