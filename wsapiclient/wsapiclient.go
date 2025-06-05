package wsapiclient

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapinet"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
	"github.com/rs/zerolog/log"
)

const (
	LibraryVersion string = "2.7.45"
)

type WSAPIClient struct {
	Caller     *httpclient.HTTPClient
	VMService  wsapivm.VMService
	NETService wsapinet.NETService
}

func New() *WSAPIClient {
	var err error
	myclient := new(WSAPIClient)
	myclient.Caller, err = httpclient.New()
	if err != nil {
		log.Error().Err(err).Msg("We can't make the Client.")
		return nil
	}
	myclient.VMService = wsapivm.New(myclient.Caller)
	myclient.NETService = wsapinet.New(myclient.Caller)
	return myclient
}
