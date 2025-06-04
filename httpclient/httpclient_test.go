package httpclient

import (
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	url, user, pass, _, insecure, debug, err := InitialData("../config.ini")
	debug = strings.ToLower(debug)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	apiClient, err := NewClient(url, user, pass, insecure, debug)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	if !strings.Contains(apiClient.BaseURL.String(), "https") || !strings.Contains(apiClient.BaseURL.String(), "http") {
		t.Errorf("The param url not contain the formatted URL: %#v", url)
	}
	if strings.Contains("none, info, error, debug", apiClient.DebugLevel) {
		t.Errorf("The Debug Level has defined a wrong level: %#v", debug)
	}
}

func TestNew(t *testing.T) {

}

func TestApiCall(t *testing.T) {

}

func TestConfigCli(t *testing.T) {

}
