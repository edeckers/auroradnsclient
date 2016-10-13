package auroradns_client

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/edeckers/auroradns_client/records"
)

func (self *AuroraDNSClient) GetRecords(zoneId string) ([]records.GetRecordsResponse, error) {
	logrus.Debugf("GetRecords(%s)", zoneId)
	relativeUrl := fmt.Sprintf("zones/%s/records", zoneId)

	response, err := self.requestor.Request(relativeUrl, "GET", []byte(""))
	if err != nil {
		logrus.Errorf("Failed to receive records: %s", err)
		return nil, err
	}

	var respData []records.GetRecordsResponse
	err = json.Unmarshal(response, &respData)
	if err != nil {
		logrus.Errorf("Failed to unmarshall response: %s", err)
		return nil, err
	}

	return respData, nil
}

func (self *AuroraDNSClient) CreateRecord(zoneId string, data records.CreateRecordRequest) (*records.CreateRecordResponse, error) {
	logrus.Debugf("CreateRecord(%s, %+v)", zoneId, data)
	body, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshall request body: %s", err)

		return nil, err
	}

	relativeUrl := fmt.Sprintf("zones/%s/records", zoneId)

	response, err := self.requestor.Request(relativeUrl, "POST", body)
	if err != nil {
		logrus.Errorf("Failed to create record: %s", err)

		return nil, err
	}

	var respData *records.CreateRecordResponse
	err = json.Unmarshal(response, &respData)
	if err != nil {
		logrus.Errorf("Failed to unmarshall response: %s", err)

		return nil, err
	}

	return respData, nil
}

func (self *AuroraDNSClient) RemoveRecord(zoneId string, recordId string) (*records.RemoveRecordResponse, error) {
	logrus.Debugf("RemoveRecord(%s, %s)", zoneId, recordId)
	relativeUrl := fmt.Sprintf("zones/%s/records/%s", zoneId, recordId)

	_, err := self.requestor.Request(relativeUrl, "DELETE", nil)
	if err != nil {
		logrus.Errorf("Failed to remove record: %s", err)

		return nil, err
	}

	return &records.RemoveRecordResponse{}, nil
}
