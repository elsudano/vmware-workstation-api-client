package wsapiclient

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapinet"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
	"github.com/rs/zerolog/log"
)

type WSAPIClient struct {
	Caller     httpclient.HTTPClient
	VMService  wsapivm.VMService
	NETService wsapinet.NETService
}

func New() (*WSAPIClient, error) {
	Caller, err := httpclient.New()
	if err != nil {
		log.Error().Err(err).Msg("We can't make the Client.")
		return nil, err
	}
	return &WSAPIClient{
		VMService:  wsapivm.New(Caller),
		NETService: wsapinet.New(Caller),
	}, nil
}
