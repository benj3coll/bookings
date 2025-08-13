package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
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
	{"gquarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"msuite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},

	{"postsearch", "/search-availability", "POST", []postData{
		{key: "start", value: "2025-09-02"},
		{key: "end", value: "2025-09-15"},
	}, http.StatusOK},
	{"postsearchjson", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2025-09-02"},
		{key: "end", value: "2025-09-15"},
	}, http.StatusOK},
	{"postreservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "j@s.com"},
		{key: "phone", value: "0049"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, el := range theTests {
		if el.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + el.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != el.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", el.name, el.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range el.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+el.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != el.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", el.name, el.expectedStatusCode, resp.StatusCode)
			}
		}

	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
