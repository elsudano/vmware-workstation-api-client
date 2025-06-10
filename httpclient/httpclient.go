package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	defaultUser       = "Admin"
	defaultPassword   = "Adm1n#00"
	defaultBaseURL    = "http://localhost:8697/api"
	defaultInsecure   = true
	defaultDebugLevel = "NONE" // DEBUG, ERROR, INFO, NONE
	// don't change this value, always activate Debug Mode
	// change behavior with ConfigCli method, it's better
	// because you can change the behavior in the future
)

// VmError that's the error that the API give us in different situations handling resources
type VmError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPClient object, this object contain:
// Client: (*http.Client) This the http client used to talk with API REST.
// BaseURL: (*url.URL) Object URL to storage URL to server.
// User: (string) Name of user to authenticate in server.
// Password: (string) Password of user, Debug: bool that show the debug it's active or not.
type HTTPClient struct {
	Client       *http.Client
	BaseURL      *url.URL
	User         string
	Password     string
	InsecureFlag bool
	DebugLevel   string
}

// NewClient constructor of the Client object
// Inputs:
// a: (string) URL address to the API REST server.
// u: (string) String with the user to connect at API REST.
// p: (string) String with the password.
// i: (bool) True if we have generated the https certificates in our API.
// d: (string) Level of Debug that we want.
// Outputs:
// *Client: (pointer) Pointer at the object Client,
// error: (error) when the client generate some error is storage in this var.
func NewClient(a string, u string, p string, i bool, d string) (*HTTPClient, error) {
	c := new(HTTPClient)
	URL, err := url.Parse(strings.TrimSpace(a))
	if err != nil {
		log.Error().Err(err).Msgf("We can't parsed the URL: %#v", err)
		return nil, err
	}
	c.BaseURL = URL
	c.User = u
	c.Password = p
	c.InsecureFlag = i
	c.DebugLevel = (strings.ToUpper(d))
	log.Debug().Msgf("Input values %#v, %#v, %#v, %#v, %#v", a, u, p, i, d)
	c.Client = &http.Client{
		Transport: &http.Transport{
			// DisableKeepAlives: false,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: i,
			},
		},
	}
	log.Debug().Msgf("Client %#v", c.Client)
	log.Info().Msg("We have created the client.")
	return c, nil
}

// New constructor of the Client object without input, this method generate a *Client
// with values by default, Return: *Client: pointer at the object Client,
// error: when the client generate some error is storage in this var.
func New() (*HTTPClient, error) {
	c, err := NewClient(defaultBaseURL, defaultUser, defaultPassword, defaultInsecure, defaultDebugLevel)
	if err != nil {
		log.Error().Err(err).Msg("We can't create the client")
		return nil, err
	}
	return c, err
}

// ConfigCli method return a pointer of Client of API but now it's configure
// Inputs:
// c: (*HTTPClient) client with all the necessary data to make a call.
// a: (string) address of URL to server of API.
// u: (string) user for to authenticate.
// p: (string) password of user.
// i: (bool) Insecure flag to http or https.
// d: (string) debug mode
func (c *HTTPClient) ConfigClient(a string, u string, p string, i bool, d string) error {
	var err error
	log.Debug().Msgf("Variables Values: %#v, %#v, %#v, %#v, %#v", a, u, p, i, d)
	c.BaseURL, err = url.Parse(a)
	if err != nil {
		log.Error().Err(err).Msg("The URL is malformed")
		return err
	}
	log.Debug().Msgf("Client BaseURL: %#v", c.BaseURL)
	c.User = u
	log.Debug().Msgf("Client User: %#v", c.User)
	c.Password = p
	log.Debug().Msgf("Client Password: %#v", c.Password)
	c.InsecureFlag = i
	log.Debug().Msgf("Client http/s: %#v", c.InsecureFlag)
	c.DebugLevel = d
	log.Debug().Msgf("Client Debug Level: %#v", c.DebugLevel)
	c.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: i,
			},
		},
	}
	log.Info().Msgf("We have configured the client.")
	return nil
}

// httpRequest method return a body of the response the API REST server,
// Input:
// c: (*HTTPClient) client with all the necessary data to make a call.
// p: (string) URL path of the API REST of the sever.
// m: (string) Type of method GET, PUT, POST, DELETE.
// pl: (bytes.Buffer) for read the Body of the request.
// Output:
// response: (io.ReadCloser) That will be the Response Body that the API give us.
// err: (error) Normal error answer that go-lang give us in a issue.
func (c *HTTPClient) ApiCall(p string, m string, pl bytes.Buffer) (io.ReadCloser, error) {
	var vmerror VmError
	req, err := http.NewRequest(m, c.RequestPath(p), &pl)
	if err != nil {
		log.Error().Err(err).Msgf("Calling to API: %#v", err)
		return nil, err
	}
	if pl.Len() > 0 {
		log.Debug().Msgf("Request Buffer: %#v", pl.String())
	}
	req.SetBasicAuth(c.User, c.Password)
	switch m {
	case "GET":
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "PUT":
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "POST":
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}
	log.Debug().Msgf("We are doing the API call")
	responseBody := new(bytes.Buffer)
	response, err := c.Client.Do(req)
	if response == nil && err != nil {
		log.Error().Err(err).Msg("The server response with timeout.")
		return nil, err
	}
	log.Debug().Msgf("Response RAW %#v", response)
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Debug().Msgf("The result of API call was: %#v", response.StatusCode)
	case http.StatusConflict:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("Trying to decode the VmError.")
			return nil, err
		}
		log.Debug().Msgf("Response StatusCode %#v Code Error %#v Message: %#v", response.StatusCode, vmerror.Code, vmerror.Message)
		// We need to create a new functionality to handle the different errors that the VmWare Workstation API give us when we have an error
		return nil, errors.New("StatusCode:" + strconv.Itoa(response.StatusCode) + ", Code Error:" + strconv.Itoa(vmerror.Code) + ", Message:" + vmerror.Message)
	case http.StatusNotFound:
		log.Debug().Msgf("Response StatusCode %#v Resource Not Found, the VmWare API isn't working properly", response.StatusCode)
		// We need to create a new functionality to handle the different errors that the VmWare Workstation API give us when we have an error
		return nil, errors.New("StatusCode:" + strconv.Itoa(response.StatusCode) + " Resource Not Found, the VmWare API isn't working properly")
	case http.StatusInternalServerError:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("Trying to decode the VmError.")
			return nil, err
		}
		log.Debug().Msgf("Response StatusCode %#v Code Error %#v Message: %#v", response.StatusCode, vmerror.Code, vmerror.Message)
		// We need to create a new functionality to handle the different errors that the VmWare Workstation API give us when we have an error
		return nil, errors.New("StatusCode:" + strconv.Itoa(response.StatusCode) + ", Code Error:" + strconv.Itoa(vmerror.Code) + ", Message:" + vmerror.Message)
	default:
		_, err = responseBody.ReadFrom(response.Body)
		if err != nil {
			log.Error().Err(err).Msgf("ResponseBody RAW %#v", responseBody)
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("The Response isn't a JSON format.")
			return nil, err
		}
		// We need to create a new functionality to handle the different errors that the VmWare Workstation API give us when we have an error
		return nil, errors.New("StatusCode:" + strconv.Itoa(response.StatusCode) + ", Code Error:" + strconv.Itoa(vmerror.Code) + ", Message:" + vmerror.Message)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error making a call")
		return nil, err
	}
	log.Debug().Msg("The API call was completed.")
	return response.Body, nil
}
