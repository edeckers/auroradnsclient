package auroradns_client

type AuroraDNSError struct {
	ErrorCode string `json:"error"`
	Message   string `json:"errormsg"`
}

func (e AuroraDNSError) Error() string {
	return e.Message
}
