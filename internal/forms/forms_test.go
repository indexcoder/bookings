package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("Form Valid Failed")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form Required Failed")
	}

	postedData := url.Values{}
	postedData.Set("a", "a")
	postedData.Set("b", "a")
	postedData.Set("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form Required Failed")
	}
}

func TestForm_Has(t *testing.T) {
	//r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("Form shows Has Failed when it does not")
	}

	// postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min length for non-existing field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Form shows min length for non-existing field")
	}

	postedData.Add("some_field", "some value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("Form shows min length for non-existing field")
	}

	postedData.Add("another_field", "abs123")
	form = New(postedData)
	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("Form shows min length for non-existing field")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows IsEmail Failed when it does not")
	}

	postedData.Add("email", "me@here.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Form shows IsEmail Failed when it does not")
	}

}
