package requests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExtendedRequest struct {
	*http.Request
}

type ExtendedHeader struct {
	http.Header
}

func (self ExtendedRequest) DebugOutput() string {
	requestBody, _ := ioutil.ReadAll(self.Request.Body)
	self.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	return fmt.Sprintf(
		"%s %s HTTP/1.1\nHost: %s\n%s\n\n%s",
		self.Request.Method, self.Request.URL.Path,
		self.Request.URL.Host,
		ExtendedHeader{self.Header}.DebugOutput(),
		requestBody)
}

func (self ExtendedHeader) DebugOutput() string {
	return fmt.Sprintf("X-AuroraDNS-Date: %s\nAuthorization: %s\nContent-Type: %s",
		self.Get("X-AuroraDNS-Date"), self.Get("Authorization"), self.Get("Content-Type"))
}
