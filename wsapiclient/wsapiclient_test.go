package wsapiclient

import "testing"

func TestNewClient(t *testing.T) {
	_, err := NewClient("", "", "", false)
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestRequestCurl(t *testing.T) {

}
