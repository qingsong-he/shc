package shc

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const contentType = "Content-Type"
const ContentTypeByJSON = "application/json"
const ContentTypeByXML = "application/xml"
const ContentTypeByFormUrlEncoded = "application/x-www-form-urlencoded"

type shc struct {
	cli *http.Client
}

func NewSHC(tripper http.RoundTripper, checkRedirect func(req *http.Request, via []*http.Request) error, jar http.CookieJar, timeout time.Duration) *shc {
	return &shc{
		cli: &http.Client{
			Transport:     tripper,
			CheckRedirect: checkRedirect,
			Jar:           jar,
			Timeout:       timeout,
		},
	}
}

func (s *shc) BuildUrlParam(bsaeUrl string, urlParam url.Values) (string, error) {
	bsaeUrlObj, err := url.Parse(bsaeUrl)
	if err != nil {
		return "", err
	}

	bsaeUrlObj.RawQuery = urlParam.Encode()
	return bsaeUrlObj.String(), nil
}

func (s *shc) BuildFormUrlEncodedReader(urlParam url.Values) io.Reader {
	return strings.NewReader(urlParam.Encode())
}

func (s *shc) BuildFormDataReader(field map[string]string, file map[string]string) (io.Reader, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for k, v := range field {
		fw, err := w.CreateFormField(k)
		if err != nil {
			return nil, "", err
		}
		fw.Write([]byte(v))
	}

	for k, v := range file {
		fileContent, err := ioutil.ReadFile(v)
		if err != nil {
			return nil, "", err
		}
		fw, err := w.CreateFormFile(k, v)
		fw.Write(fileContent)
	}
	w.Close()

	return &b, w.FormDataContentType(), nil
}

func (s *shc) Request(method string, url string, contentTypeValue string, header http.Header, body io.Reader) (*http.Request, *http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}

	if contentTypeValue != "" {
		req.Header.Add(contentType, contentTypeValue)
	}

	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return req, resp, nil
}
