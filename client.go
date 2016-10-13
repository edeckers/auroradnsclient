package auroradns_client

import (
	"github.com/edeckers/auroradns_client/requests"
)

type AuroraDNSClient struct {
	requestor *requests.AuroraRequestor
}

func NewAuroraDNSClient(baseUrl string, userId string, key string) (*AuroraDNSClient, error) {
	requestor, err := requests.NewAuroraRequestor(baseUrl, userId, key)
	if err != nil {
		return nil, err
	}

	return &AuroraDNSClient{
		requestor: requestor,
	}, nil
}
