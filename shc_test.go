package shc

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestShc_BuildUrlParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r.FormValue("a"))
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	shc := NewSHC(nil, nil, nil, 0)

	urlByParam, err := shc.BuildUrlParam(ts.URL, url.Values{"a": []string{"b"}})
	if err != nil {
		t.Fatal(err)
	}

	req, resp, err := shc.Request(http.MethodGet, urlByParam, "", http.Header{"c": []string{"d"}}, nil)
	if err != nil {
		t.Fatal(err)
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(reqDump))

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(respDump))
}

func TestShc_BuildFormDataReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(0xffff)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r.FormValue("a"))
		file, fileHeader, err := r.FormFile("bashrc")
		if err != nil {
			t.Fatal(err)
		}
		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		t.Log(fileHeader.Filename, fileHeader.Size, fileHeader.Header)
		t.Log(string(fileContent))
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	shc := NewSHC(nil, nil, nil, 0)

	formDataReader, contentTypeByFormData, err := shc.BuildFormDataReader(map[string]string{"a": "b"}, map[string]string{"bashrc": "/etc/bashrc"})
	if err != nil {
		t.Fatal(err)
	}

	req, resp, err := shc.Request(http.MethodPost, ts.URL, contentTypeByFormData, http.Header{"c": []string{"d"}}, formDataReader)
	if err != nil {
		t.Fatal(err)
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(reqDump))

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(respDump))
}

func TestShc_BuildFormUrlEncodedReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r.FormValue("c"))
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	shc := NewSHC(nil, nil, nil, 0)

	req, resp, err := shc.Request(http.MethodPost, ts.URL, ContentTypeByFormUrlEncoded, http.Header{"a": []string{"b"}}, shc.BuildFormUrlEncodedReader(url.Values{"c": []string{"d"}}))
	if err != nil {
		t.Fatal(err)
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(reqDump))

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(respDump))
}
