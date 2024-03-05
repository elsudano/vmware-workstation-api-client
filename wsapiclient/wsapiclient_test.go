package wsapiclient

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	url, user, pass, _, insecure, debug, err := InitialData("../config.ini")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	apiClient, err := NewClient(url, user, pass, insecure, debug)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	VM, err := apiClient.GetVM("545OMDAL1R520604HKNKA6TTK6TBNOHK")
	if VM.Denomination != "parentvm" || err != nil {
		t.Errorf("[ERROR][WSAPICLI] Fi: wsapiclient_test.go Fu: TestNewClient M: You need make sure that the ParentVM it's called 'parentvm' %#v", err)
	}
}
