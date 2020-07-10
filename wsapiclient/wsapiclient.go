package wsapiclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	libraryVersion  = "0.1.0"
	defaultUser     = "admin"
	defaultPassword = "Adm1n#00"
	defaultBaseURL  = "https://localhost:8697/api"
	defaultInsecure = true
	defaultDebug    = true
	// don't change this value, always activate Debug Mode
	// change behavior with ConfigCli method, it's better
	// because you can change the behavior in the future
)

// Client object, this object contain: Client: *http.Client this the http client used to talk with API REST,
// BaseURL: *url.URL object URL to storage URL to server, User: string name of user to authenticate in server
// Password: string password of user, Debug: bool that show the debug it's active or not
type Client struct {
	// HTTP client used to communicate with the DO API.
	Client *http.Client
	// Base URL for API requests.
	BaseURL *url.URL
	// User to access
	User string
	// Password of User
	Password string
	// Insecure Mode
	InsecureFlag bool
	// Debug Mode
	Debug bool
}

// NewClient constructor of the Client object Input: a: URL address to the API REST server
// u: string with the user to connect at API REST, p: string with the password,
// d: bool to activate or not the debug, Return: *Client: pointer at the object Client,
// error: when the client generate some error is storage in this var.
func NewClient(a string, u string, p string, i bool, d bool) (*Client, error) {
	c := new(Client)
	c.BaseURL, _ = url.Parse(a)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: NewClient Obj:URL %#v\n", c.BaseURL)
	c.User = u
	c.Password = p
	c.Debug = d
	// Como estamos desarrollando una API que se encaga de comunicarnos con
	// VmWare Workstation Pro API REST, y la propia API de Workstation se
	// accede a través de un cliente web necesitamos generar un cliente
	// para poder realizar las peticiones.
	//
	// Por eso con esta variable creamos un cliente el cual nos permite realizar
	// peticiones sobre paginas https pero sin la necesidad de la verificaión
	// del certificado de la pagina, esto se hace por si la pagina tiene un
	// certificado auto-firmado.
	//
	// Esto solo es valido para testing, se ha de quitar para validar correctamente
	c.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: i,
			},
		},
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: NewClient Obj:web Client %#v\n", c.Client)
	return c, nil
}

// New constructor of the Client object without input, this method generate a *Client
// with values by default, Return: *Client: pointer at the object Client,
// error: when the client generate some error is storage in this var.
func New() (*Client, error) {
	c, err := NewClient(defaultBaseURL, defaultUser, defaultPassword, defaultInsecure, defaultDebug)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: New Obj:api Client %#v\n", c)
	return c, err
}

// Method of object *Client to change the debug parameter
// activate or disable the debug mode of the *Client.
func (c *Client) SwitchDebug() {
	// for config Debug mode
	if c.Debug {
		log.SetOutput(ioutil.Discard)
		c.Debug = false
	}
	if !c.Debug {
		log.SetOutput(os.Stdout)
		c.Debug = true
	}
}

// Method ConfigCli return a pointer of Client of API but now it's configure
// Inputs: a: address of URL to server of API u: user for to authenticate
// p: password of user, i: Insecure flag to http or https, d: debug mode
func (c *Client) ConfigCli(a string, u string, p string, i bool, d bool) {
	var err error
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:Variables %#v, %#v, %#v, %#v\n", a, u, p, d)
	// for config Debug mode
	if !d {
		log.SetOutput(ioutil.Discard)
		c.Debug = false
	}
	c.BaseURL, err = url.Parse(a)
	if err != nil {
		panic(err)
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:%#v\n", c.BaseURL)
	c.User = u
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:%#v\n", c.User)
	c.Password = p
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:%#v\n", c.Password)
	c.InsecureFlag = i
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:%#v\n", c.InsecureFlag)
	c.Debug = d
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Obj:%#v\n", c.Debug)

}

// Method httpRequest return a body of the response the API REST server, Input:
// p: URL path of the API REST of the sever, m: Type of method GET, PUT, POST, DELETE
// pl: bytes.Buffer for read the Body of the request, Return: cl:
func (c *Client) httpRequest(p string, m string, pl bytes.Buffer) (io.ReadCloser, error) {
	req, err := http.NewRequest(m, c.requestPath(p), &pl)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapiclient.go Fu: httpRequest Obj:request error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: httpRequest Obj:buffer %#v\n", pl.String())
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
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: httpRequest Obj:request %#v\n", req)
	response, err := c.Client.Do(req)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapiclient.go Fu: httpRequest Obj:response error %#v\n", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		responseBody := new(bytes.Buffer)
		_, err := responseBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("[WSAPICLI][ERROR] Fi: wsapiclient.go Fu: httpRequest Obj:respBody %#v\n", responseBody)
			return nil, err
		}
		log.Printf("[WSAPICLI][ERROR] Fi: wsapiclient.go Fu: httpRequest Obj:StatusCode %#v Body %#v\n", response.StatusCode, responseBody.String())
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: httpRequest Obj:response %#v\n", response)
	return response.Body, nil
}

// Method requestPath show the URL to the request of httpClient, Input:
// p: string just the path of the URL, Return: string with the complete URL to access
func (c *Client) requestPath(p string) string {
	r := fmt.Sprintf("%s/%s", c.BaseURL, p)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: requestPath Obj:%#v\n", r)
	return r
}
