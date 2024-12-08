package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key string
	val string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"features", "/features", "GET", []postData{}, http.StatusOK},
	{"about-summary", "/about-summary", "GET", []postData{}, http.StatusOK},
	{"search", "/search", "POST", []postData{
		{key: "start", val: "2021-01-01"},
		{key: "end", val: "2021-01-02"},
	}, http.StatusOK},
	{"search-json", "/search-json", "POST", []postData{
		{key: "start", val: "2021-01-01"},
		{key: "end", val: "2021-01-02"},
	}, http.StatusOK},
	{"about-post", "/about", "POST", []postData{
		{key: "start", val: "2021-01-01"},
		{key: "end", val: "2021-01-02"},
	}, http.StatusOK},
}

func TestNewHandler(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, tc := range theTests {
		if tc.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + tc.url)
			if err != nil {
				t.Log(err.Error())
				t.Fail()
			}
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, tc.expectedStatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range tc.params {
				values.Add(x.key, x.val)
			}
			resp, err := ts.Client().PostForm(ts.URL+tc.url, values)
			if err != nil {
				t.Log(err.Error())
				t.Fail()
			}
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, tc.expectedStatusCode)
			}
		}
	}
}
