package auroradns_client

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/edeckers/auroradns_client/zones"
)

func (self *AuroraDNSClient) GetZones() ([]zones.ZoneRecord, error) {
	logrus.Debugf("GetZones")
	response, err := self.requestor.Request("zones", "GET", []byte(""))

	if err != nil {
		logrus.Errorf("Failed to get zones: %s", err)
		return nil, err
	}

	var respData []zones.ZoneRecord
	err = json.Unmarshal(response, &respData)
	if err != nil {
		logrus.Errorf("Failed to unmarshall response: %s", err)
		return nil, err
	}

	logrus.Debugf("Unmarshalled response: %+v", respData)
	return respData, nil
}
