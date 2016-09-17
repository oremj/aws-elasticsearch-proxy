package proxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type ElasticSearch struct {
	Endpoint string
	Region   string
	Signer   *v4.Signer

	bufferPool *BufferPool
}

func NewElasticSearch(endpoint, region string, signer *v4.Signer) *ElasticSearch {
	return &ElasticSearch{
		Endpoint: endpoint,
		Region:   region,
		Signer:   signer,

		bufferPool: NewBufferPool(),
	}
}

// req.Body will be closed after this.
func (e *ElasticSearch) buildRequest(req *http.Request) (*http.Request, error) {
	var body *bytes.Buffer
	if req.Body != nil {
		body = new(bytes.Buffer)
		e.copy(body, req.Body)
		req.Body.Close()
	}

	newreq := &http.Request{
		Method: req.Method,
		URL: &url.URL{
			Scheme:   "https",
			Host:     e.Endpoint,
			Path:     req.URL.Path,
			RawQuery: req.URL.RawQuery,
		},
		Proto:      "HTTP/1.1",
		ProtoMinor: 1,
		ProtoMajor: 1,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(body),
		Host:       e.Endpoint,
	}
	if body != nil {
		newreq.ContentLength = int64(body.Len())
	}

	copyHeader(newreq.Header, req.Header)
	for _, h := range hopHeaders {
		newreq.Header.Del(h)
	}

	_, err := e.Signer.Sign(newreq, bytes.NewReader(body.Bytes()), "es", e.Region, time.Now())
	if err != nil {
		return nil, fmt.Errorf("Error Signing Request: %v", err)
	}
	return newreq, nil
}

func (e *ElasticSearch) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newreq, err := e.buildRequest(req)
	if err != nil {
		log.Println("buildRequest: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(newreq)
	if err != nil {
		log.Printf("making request: %v", err)
		http.Error(w, "AWS Error", http.StatusBadGateway)
		return
	}

	for _, h := range hopHeaders {
		resp.Header.Del(h)
	}
	copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)

	if resp.Body != nil {
		e.copy(w, resp.Body)
	}
}

func (e *ElasticSearch) copy(dst io.Writer, src io.Reader) {
	buf := e.bufferPool.Get()
	io.CopyBuffer(dst, src, buf)
	e.bufferPool.Put(buf)
}

// Taken from httputil/reverseproxy.go
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
