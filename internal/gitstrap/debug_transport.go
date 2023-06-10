package gitstrap

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type logTransport struct {
	origin http.RoundTripper
	tag    string
}

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("[%s] >>> %s %s", t.tag, req.Method, req.URL)
	if req.Body != nil {
		defer req.Body.Close()
		if data, err := io.ReadAll(req.Body); err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(data))
			log.Print(string(data))
		}
	}
	rsp, err := t.origin.RoundTrip(req)
	if err != nil {
		log.Printf("[%s] %s ERR: %s", t.tag, req.URL, err)
	} else {
		log.Printf("[%s] %s <<< %d", t.tag, req.URL, rsp.StatusCode)
		if rsp.Body != nil {
			defer rsp.Body.Close()
			if data, err := io.ReadAll(rsp.Body); err == nil {
				rsp.Body = io.NopCloser(bytes.NewBuffer(data))
				log.Print(string(data))
			}
		}
	}
	return rsp, err
}
