package wsapiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	libraryVersion    = "1.1.17"
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

// Client object, this object contain: Client: *http.Client this the http client used to talk with API REST,
// BaseURL: *url.URL object URL to storage URL to server, User: string name of user to authenticate in server
// Password: string password of user, Debug: bool that show the debug it's active or not
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
	c.BaseURL, _ = url.Parse(a)
	c.User = u
	c.Password = p
	c.InsecureFlag = i
	c.DebugLevel = (strings.ToUpper(d))
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: NewClient Input values %#v, %#v, %#v, %#v, %#v\n", a, u, p, i, d)
	}
	c.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: i,
			},
		},
	}
	if c.DebugLevel != "NONE" {
		log.SetOutput(os.Stderr)
	}
	if c.DebugLevel == "NONE" {
		log.SetOutput(io.Discard)
	}
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: NewClient  Client %#v\n", c.Client)
	}
	if c.DebugLevel == "INFO" {
		log.Printf("[INFO][WSAPICLI] Fi: wsapiclient.go Fu: NewClient we completed tasks")
	}
	return c, nil
}

// New constructor of the Client object without input, this method generate a *Client
// with values by default, Return: *Client: pointer at the object Client,
// error: when the client generate some error is storage in this var.
func New() (*Client, error) {
	c, err := NewClient(defaultBaseURL, defaultUser, defaultPassword, defaultInsecure, defaultDebugLevel)
	if err != nil && c.DebugLevel == "ERROR" {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: New Error creating the client %#v\n", err)
	}
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: New api Client %#v\n", c)
	}
	if c.DebugLevel == "INFO" {
		log.Printf("[INFO][WSAPICLI] Fi: wsapiclient.go Fu: New we completed tasks")
	}
	return c, err
}

// SwitchDebug method of object *Client to change the debug parameter
// activate or disable the debug mode of the *Client.
// Imputs:
// c: (pointer) The pointer of the client that we are using
// l: (string) The Level of the Debug that we want.
func (c *Client) SwitchDebugLevel(l string) {
	switch l {
	case "NONE":
		log.SetOutput(io.Discard)
		c.DebugLevel = "NONE"
	case "INFO":
		log.SetOutput(os.Stderr)
		c.DebugLevel = "INFO"
	case "ERROR":
		log.SetOutput(os.Stderr)
		c.DebugLevel = "ERROR"
	case "DEBUG":
		log.SetOutput(os.Stderr)
		c.DebugLevel = "DEBUG"
	default:
		log.SetOutput(io.Discard)
		c.DebugLevel = "NONE"
	}
	if c.DebugLevel == "INFO" {
		log.Printf("[INFO][WSAPICLI] Fi: wsapiclient.go Fu: SwitchDebugLevel we completed tasks")
	}
}

// ConfigCli method return a pointer of Client of API but now it's configure
// Inputs: a: address of URL to server of API u: user for to authenticate
// p: password of user, i: Insecure flag to http or https, d: debug mode
func (c *Client) ConfigCli(a string, u string, p string, i bool, d string) {
	var err error
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Variables %#v, %#v, %#v, %#v, %#v\n", a, u, p, i, d)
	}
	c.BaseURL, err = url.Parse(a)
	if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Api Client %#v\n", err)
	}
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli BaseURL: %#v\n", c.BaseURL)
	}
	c.User = u
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli User: %#v\n", c.User)
	}
	c.Password = p
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Password: %#v\n", c.Password)
	}
	c.InsecureFlag = i
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli http/s: %#v\n", c.InsecureFlag)
	}
	c.DebugLevel = d
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Debug Level: %#v\n", c.DebugLevel)
	}
	c.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: i,
			},
		},
	}
	if c.DebugLevel != "NONE" {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}
	if c.DebugLevel == "INFO" {
		log.Printf("[INFO][WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli we completed tasks")
	}
}

// httpRequest method return a body of the response the API REST server, Input:
// p: URL path of the API REST of the sever, m: Type of method GET, PUT, POST, DELETE
// pl: bytes.Buffer for read the Body of the request, Return: cl:
func (c *Client) httpRequest(p string, m string, pl bytes.Buffer) (io.ReadCloser, VmError, error) {
	var vmerror VmError
	req, err := http.NewRequest(m, c.requestPath(p), &pl)
	if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest request error %#v\n", err)
		return nil, vmerror, err
	}
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Buffer %#v\n", pl.String())
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
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Request before that run %#v\n", req)
	}
	// in this line we will need to create a management of queue
	responseBody := new(bytes.Buffer)
	response, err := c.Client.Do(req)
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if c.DebugLevel == "DEBUG" {
			log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  StatusCode %#v\n", response.StatusCode)
		}
	case http.StatusConflict:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Decoding Response Body %#v\n", err)
			return nil, vmerror, err
		}
		if c.DebugLevel == "DEBUG" {
			log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  StatusCode %#v Code Error %#v Message: %#v\n", response.StatusCode, vmerror.Code, vmerror.Message)
		}
		return nil, vmerror, err
	case http.StatusInternalServerError:
		err = json.NewDecoder(response.Body).Decode(&vmerror)
		if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Decoding Response Body %#v\n", err)
			return nil, vmerror, err
		}
		if c.DebugLevel == "DEBUG" {
			log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  StatusCode %#v Code Error %#v Message: %#v\n", response.StatusCode, vmerror.Code, vmerror.Message)
		}
		return nil, vmerror, err
	default:
		_, err = responseBody.ReadFrom(response.Body)
		if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  ResponseBody RAW %#v\n", responseBody)
			return nil, vmerror, err
		}
		if c.DebugLevel == "DEBUG" {
			log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Response Body before %#v\n", responseBody.String())
		}
		err = json.NewDecoder(responseBody).Decode(&vmerror)
		if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest Message: I can't read the json structure %s", err)
			return nil, vmerror, err
		}
		if c.DebugLevel == "DEBUG" {
			log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Response RAW %#v\n", response)
		}
		return nil, vmerror, err
	}
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Response after run %#v\n", response)
	}
	if err != nil && (c.DebugLevel == "ERROR" || c.DebugLevel == "DEBUG") {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapiclient.go Fu: httpRequest  Response error %#v\n", err)
		return nil, vmerror, err
	}
	return response.Body, vmerror, err
}

// requestPath method show the URL to the request of httpClient.
// Input:
// p: string just the path of the URL.
// Return:
// string with the complete URL to access
func (c *Client) requestPath(p string) string {
	r := fmt.Sprintf("%s/%s", c.BaseURL, p)
	if c.DebugLevel == "DEBUG" {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapiclient.go Fu: requestPath %#v\n", r)
	}
	return r
}
