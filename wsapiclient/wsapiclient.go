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
	libraryVersion  = "0.0.1"
	defaultUser     = "admin"
	defaultPassword = "Adm1n#00"
	defaultBaseURL  = "https://localhost:8697/api"
	defaultDebug    = true
)

type Client struct {
	// HTTP client used to communicate with the DO API.
	Client *http.Client
	// Base URL for API requests.
	BaseURL *url.URL
	// User to access
	User string
	// Password of User
	Password string
	// Debug Mode
	Debug bool
}

func NewClient(a string, u string, p string, d bool) (*Client, error) {
	c := new(Client)
	// for config Debug mode
	if !d {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}
	c.BaseURL, _ = url.Parse(a)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: NewClient Ob: %#v\n", c.BaseURL)
	c.User = u
	c.Password = p
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
				InsecureSkipVerify: true,
			},
		},
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: NewClient Ob: %#v\n", c.Client)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: NewClient Ob: %#v\n", c)
	return c, nil
}

func New() (*Client, error) {
	c, err := NewClient(defaultBaseURL, defaultUser, defaultPassword, defaultDebug)
	log.SetOutput(os.Stdout)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: New Ob: %#v\n", c)
	log.SetOutput(ioutil.Discard)
	return c, err
}

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

func (c *Client) ConfigCli(a string, u string, p string, d bool) {
	var err error
	// for config Debug mode
	if !d {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}
	c.BaseURL, err = url.Parse(a)
	if err != nil {
		panic(err)
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Ob: %#v\n", c.BaseURL)
	c.User = u
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Ob: %#v\n", c.User)
	c.Password = p
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: ConfigCli Ob: %#v\n", c.Password)

}

func (c *Client) httpRequest(p string, m string, pl bytes.Buffer) (cl io.ReadCloser, err error) {
	req, err := http.NewRequest(m, c.requestPath(p), &pl)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.User, c.Password)
	switch m {
	case "GET":
	case "PUT":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: httpRequest Ob: %#v\n", resp)
	return resp.Body, nil
}

func (c *Client) requestPath(p string) string {
	r := fmt.Sprintf("%s/%s", c.BaseURL, p)
	log.Printf("[WSAPICLI] Fi: wsapiclient.go Fu: requestPath Ob: %#v\n", r)
	return r
}
