package wsapiclient

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapinet"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
	"github.com/rs/zerolog/log"
)

const (
	libraryVersion string = "2.7.45"
)

type WSAPIClient struct {
	VMService  wsapivm.VMService
	NETService wsapinet.NETService
}

func New() *WSAPIClient {
	Caller, err := httpclient.New()
	if err != nil {
		log.Error().Err(err).Msg("We can't make the Client.")
		return nil
	}
	return &WSAPIClient{
		VMService:  wsapivm.New(Caller),
		NETService: wsapinet.New(Caller),
	}
}

// func ConfigApiClient(a string, u string, p string, i bool, d string) error {
// 	return &WSAPIClient.VMService.ConfigClient(a, u, p, i, d)
// }
