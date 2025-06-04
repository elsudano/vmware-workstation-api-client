package httpclient

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
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
