package wsapiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	libraryVersion    = "1.2.20"
	defaultUser       = "Admin"
	defaultPassword   = "Adm1n#00"
	defaultBaseURL    = "http://localhost:8697/api"
	defaultInsecure   = true
	defaultDebugLevel = "NONE" // DEBUG, ERROR, INFO, NONE
	// don't change this value, always activate Debug Mode
	// change behavior with ConfigCli method, it's better
	// because you can change the behavior in the future
)

// VmError that's rhe error that the API give us in different situations handling resources
type VmError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Client object, this object contain:
// Client: (*http.Client) This the http client used to talk with API REST.
// BaseURL: (*url.URL) Object URL to storage URL to server.
// User: (string) Name of user to authenticate in server.
// Password: (string) Password of user, Debug: bool that show the debug it's active or not.
type Client struct {
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
func NewClient(a string, u string, p string, i bool, d string) (*Client, error) {
	c := new(Client)
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
	ConfigLog(c.DebugLevel, "HR")
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
func New() (*Client, error) {
	c, err := NewClient(defaultBaseURL, defaultUser, defaultPassword, defaultInsecure, defaultDebugLevel)
	log.Debug().Msgf("Client Object %#v", c)
	log.Error().Err(err).Msg("We can't create the client")
	log.Info().Msg("We have created the client.")
	return c, err
}

// SwitchDebug method of object *Client to change the debug parameter
// activate or disable the debug mode of the *Client.
// Imputs:
// c: (pointer) The pointer of the client that we are using
// l: (string) The Level of the Debug that we want.
func (c *Client) SwitchDebugLevel(l string) {
	switch l {
	case "INFO":
		c.DebugLevel = "INFO"
	case "ERROR":
		c.DebugLevel = "ERROR"
	case "DEBUG":
		c.DebugLevel = "DEBUG"
	default:
		c.DebugLevel = "NONE"
	}
	log.Debug().Str("level", c.DebugLevel).Msg("We have changed the Log Level at: ")
	log.Info().Msg("We have changed the Log Level.")
}

// ConfigCli method return a pointer of Client of API but now it's configure
// Inputs: a: address of URL to server of API u: user for to authenticate
// p: password of user, i: Insecure flag to http or https, d: debug mode
func (c *Client) ConfigCli(a string, u string, p string, i bool, d string) {
	var err error
	log.Debug().Msgf("Variables Values: %#v, %#v, %#v, %#v, %#v", a, u, p, i, d)
	c.BaseURL, err = url.Parse(a)
	log.Error().Err(err).Msg("The URL is malformed")
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
}

// ConfigLog method change the behavior that how to handle the logging on our API
// Inputs:
// l: (strings) Which will be the level bu default that we want in console
// f: (string) Which will be the format that we want in the console, JSON or HR (Human Readable)
func ConfigLog(l string, f string) {
	l = strings.ToUpper(l)
	switch l {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	// Global Settings https://github.com/rs/zerolog?tab=readme-ov-file#global-settings
	zerolog.TimeFieldFormat = time.RFC3339 // zerolog.TimeFormatUnix zerolog.TimeFormatUnixMs, zerolog.TimeFormatUnixMicro

	// Customized Fields Name https://github.com/rs/zerolog?tab=readme-ov-file#customize-automatic-field-names
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	// To trace the errors https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()

	// Formatting https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
	switch strings.ToUpper(f) {
	case "CONSOLE":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		ConsoleWriter.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "HR":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-6s][WSAPICLI]", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	default:
		log.Logger = log.With().Logger().Output(os.Stderr)
	}
}

// httpRequest method return a body of the response the API REST server, Input:
// p: URL path of the API REST of the sever, m: Type of method GET, PUT, POST, DELETE
// pl: bytes.Buffer for read the Body of the request, Return: cl:
func (c *Client) httpRequest(p string, m string, pl bytes.Buffer) (io.ReadCloser, VmError, error) {
	var vmerror VmError
	req, err := http.NewRequest(m, c.requestPath(p), &pl)
	if err != nil {
		log.Error().Err(err).Msgf("Calling to API: %#v", err)
		return nil, vmerror, err
	}
	if pl.Len() > 0 {
		log.Debug().Msgf("Request Buffer: %#v", pl.String())
	}
	req.SetBasicAuth(c.User, c.Password)
	switch m {
	case "GET":
		// req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "PUT":
		// req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "POST":
		// req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
		req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}
	log.Debug().Msgf("We are doing the API call")
	// in this line we will need to create a management of queue
	responseBody := new(bytes.Buffer)
	response, err := c.Client.Do(req)
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Debug().Msgf("The result of API call was: %#v", response.StatusCode)
	case http.StatusConflict:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("Trying decode the answer.")
			return nil, vmerror, err
		}
		log.Debug().Msgf("Response StatusCode %#v Code Error %#v Message: %#v", response.StatusCode, vmerror.Code, vmerror.Message)
		return nil, vmerror, err
	case http.StatusInternalServerError:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("Trying decode the Response.")
			return nil, vmerror, err
		}
		log.Debug().Msgf("Response StatusCode %#v Code Error %#v Message: %#v", response.StatusCode, vmerror.Code, vmerror.Message)
		return nil, vmerror, err
	default:
		_, err = responseBody.ReadFrom(response.Body)
		if err != nil {
			log.Error().Err(err).Msgf("ResponseBody RAW %#v", responseBody)
			return nil, vmerror, err
		}
		err = json.NewDecoder(responseBody).Decode(&vmerror)
		if err != nil {
			log.Error().Err(err).Msg("The Response isn't a JSON format.")
			return nil, vmerror, err
		}
		return nil, vmerror, err
	}
	log.Debug().Msgf("Response RAW %#v", response)
	if err != nil {
		log.Error().Err(err).Msg("Response error")
		return nil, vmerror, err
	}
	log.Debug().Msg("The API call was completed successfully.")
	return response.Body, vmerror, err
}

// requestPath method show the URL to the request of httpClient.
// Input:
// p: (string) just the path of the URL.
// Return:
// (string) with the complete URL to access
func (c *Client) requestPath(p string) string {
	r := fmt.Sprintf("%s/%s", c.BaseURL, p)
	log.Debug().Str("URL", c.BaseURL.Host+"/"+p).Msg("The whole endpoint that we will visit.")
	return r
}
