package requests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	request_errors "github.com/edeckers/auroradns_client/requests/errors"
	"github.com/edeckers/auroradns_client/tokens"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type AuroraRequestor struct {
	endpoint string
	userId   string
	key      string
}

func NewAuroraRequestor(endpoint string, userId string, key string) (*AuroraRequestor, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("Aurora endpoint missing")
	}

	if userId == "" || key == "" {
		return nil, fmt.Errorf("Aurora credentials missing")
	}

	return &AuroraRequestor{endpoint: endpoint, userId: userId, key: key}, nil
}

func (self *AuroraRequestor) buildRequest(relativeUrl string, method string, body []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", self.endpoint, relativeUrl)

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		logrus.Errorf("Failed to build request: %s", err)

		return request, err
	}

	timestamp := time.Now().UTC()
	fmtTime := timestamp.Format("20060102T150405Z")

	token := tokens.NewToken(self.userId, self.key, method, fmt.Sprintf("/%s", relativeUrl), timestamp)

	request.Header.Set("X-AuroraDNS-Date", fmtTime)
	request.Header.Set("Authorization", fmt.Sprintf("AuroraDNSv1 %s", token))

	request.Header.Set("Content-Type", "application/json")

	rawRequest, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		logrus.Errorf("Failed to dump request: %s", err)
	}

	logrus.Debugf("Built request:\n%s", rawRequest)

	return request, err
}

func (self *AuroraRequestor) testInvalidResponse(resp *http.Response, response []byte) ([]byte, error) {
	if resp.StatusCode < 400 {
		return response, nil
	}

	logrus.Errorf("Received invalid status code %d:\n%s", resp.StatusCode, response)

	content := errors.New(string(response))

	statusCodeErrorMap := map[int]error{
		400: request_errors.BadRequest(content),
		401: request_errors.Unauthorized(content),
		403: request_errors.Forbidden(content),
		404: request_errors.NotFound(content),
		500: request_errors.ServerError(content),
	}

	mappedError := statusCodeErrorMap[resp.StatusCode]

	if mappedError == nil {
		return nil, request_errors.InvalidStatusCodeError(content)
	}

	return nil, mappedError
}

func (self *AuroraRequestor) Request(relativeUrl string, method string, body []byte) ([]byte, error) {
	req, err := self.buildRequest(relativeUrl, method, body)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed request: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	rawResponse, err := httputil.DumpResponse(resp, true)
	logrus.Debugf("Received raw response:\n%s", rawResponse)
	if err != nil {
		logrus.Errorf("Failed to dump response: %s", err)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response: %s", response)
		return nil, err
	}

	response, err = self.testInvalidResponse(resp, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
