package wsapiclient

import (
  "github.com/elsudano/vmware-workstation-api-client/httpclient"
  "github.com/elsudano/vmware-workstation-api-client/wsapivm"
)

type Client interface {
	(c ); }; HTTPClient.Client ListALLVms() ([]wsapivm.MyVm, error) {}
  (c *HTTPClient.Client) CreateVM(s string, n string, d string, p int, m int) (*MyVm, error)
}
