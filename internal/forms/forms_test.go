package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	f := New(url.Values{})

	if len(f.Errors) != 0 {
		t.Error("error in NEW")
	}
}

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid when we should have got valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows not having required fields when it has")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.Has("whatever") {
		t.Error("form shows it has field when does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "b")
	form = New(postedData)
	if !form.Has("a") {
		t.Error("form says it does not have field a when it has")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.MinLength("a", 1) {
		t.Error("form shows min length for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "bbb")
	form = New(postedData)
	if !form.MinLength("a", 2) {
		t.Error("form shows field a with length 3 does not meet min length 2")
	}

	if form.MinLength("a", 99) {
		t.Error("form shows field a meets min length when it does not")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "bbb")
	form = New(postedData)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("got a valid email for existenting field with no email")
	}

	postedData = url.Values{}
	postedData.Add("a", "b@b.c")
	form = New(postedData)
	form.IsEmail("a")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

}
