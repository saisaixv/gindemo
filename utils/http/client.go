package http

import (
	"bytes"
	"crypto/tls"
	"net"
	// "crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// NOTE: inherit from http.DefaultTransport, and skip server cert verify
var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}

type readCloser struct {
	io.Reader
}

func (r readCloser) Close() error {
	return nil
}

func CallAPI(method, url string, content []byte, header http.Header, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = header
	}
	req.ContentLength = int64(len(content))
	if req.ContentLength > 0 {
		req.Body = readCloser{bytes.NewReader(content)}
	}

	client := &http.Client{Transport: defaultTransport}
	// client := new(http.Client)
	client.Timeout = timeout
	return client.Do(req)
}

func CallJSONAPI(method, url string, data interface{}, header http.Header, timeout time.Duration) (*http.Response, error) {
	var content []byte
	var err error
	if data != nil {
		content, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if header == nil {
			header = http.Header{}
		}
		header.Set("ContentType", "application/json")
	}
	return CallAPI(method, url, content, header, timeout)
}

func ResponseBody(rsp *http.Response) ([]byte, error) {
	return FetchBody(rsp.Body)
}

func FetchBody(rc io.ReadCloser) ([]byte, error) {
	defer rc.Close()
	return ioutil.ReadAll(rc)
}
