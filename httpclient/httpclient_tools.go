package httpclient

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// SwitchDebug method of object *Client to change the debug parameter
// activate or disable the debug mode of the *Client.
// Imputs:
// c: (pointer) The pointer of the client that we are using
// l: (string) The Level of the Debug that we want.
func (c *HTTPClient) SwitchDebugLevel(l string) {
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

// requestPath method show the URL to the request of httpClient.
// Input:
// p: (string) just the path of the URL.
// Return:
// (string) with the complete URL to access
func (c *HTTPClient) RequestPath(p string) string {
	r := fmt.Sprintf("%s/%s", c.BaseURL, p)
	log.Debug().Str("URL", c.BaseURL.Host+"/"+p).Msg("The whole endpoint that we will visit.")
	return r
}

// InitialData is a extra function to fill the fields that we need to use
// to make the tests in our API.
// Inputs:
// f: (string) is the file where we have the configuration.
// Outputs:
// url: (string) That will be the URL of our API endpoint
// user: (string) That will be the User of our API
// pass: (string) That will be the Password of our API
// parentid: (string) That will be the Parent ID of our VM
// insecure: (bool) If our API works with HTTP we will set true here
// debug: (bool) If we need troubleshot our API we will set true here
// error: (error) When the function catch some error print it here
func InitialData(f string) (string, string, string, string, bool, string, error) {
	var user, pass, url, debug, parentid string
	var insecure = false
	fileInfo, err := os.Stat(f)
	if err != nil {
		log.Error().Err(err).Msg("While we trying to check the file.")
		return "", "", "", "", false, "", err
	}
	if fileInfo.Mode().IsDir() {
		log.Error().Err(err).Msg("It is a directory, please select a config file.")
		return "", "", "", "", false, "", nil
	} else if fileInfo.Mode().IsRegular() {
		file, err := os.Open(f)
		if err != nil {
			log.Error().Err(err).Msg("Failed opening file, please make sure the config file exists.")
			return "", "", "", "", false, "", err
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			temp := strings.SplitN(scanner.Text(), ":", 2)
			key := strings.ToLower(temp[0])
			if key == "user" {
				user = strings.TrimSpace(temp[1])
			}
			if key == "password" {
				pass = strings.TrimSpace(temp[1])
			}
			if key == "baseurl" {
				url = strings.TrimSpace(temp[1])
			}
			if key == "parentid" {
				parentid = strings.TrimSpace(temp[1])
			}
			if key == "insecure" && strings.TrimSpace(temp[1]) == "true" {
				insecure = true
			}
			if key == "debug" {
				debug = strings.TrimSpace(temp[1])
			}
		}
	} else {
		log.Error().Msgf("We haven't handled this error, something was wrong, please try again")
		return "", "", "", "", false, "", nil
	}
	return url, user, pass, parentid, insecure, debug, nil
}
