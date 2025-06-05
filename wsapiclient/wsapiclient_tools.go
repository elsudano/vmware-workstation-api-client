package wsapiclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigLog method change the behavior that how to handle the logging on our API
// Inputs:
// l: (strings) Which will be the level bu default that we want in console
// f: (string) The format we want:
// (defaut) JSON: in stderr,
// FILE json format in debug.log file,
// CONSOLE stdout with color in json format,
// ERROR stderr in json format
// HR (Human Readable) in stdout
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
	case "FILE":
		file, err := os.OpenFile(
			"debug.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		ConsoleWriter := zerolog.ConsoleWriter{Out: file, NoColor: false}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
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
		ConsoleWriter.PartsExclude = []string{
			zerolog.TimestampFieldName,
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "ERROR":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-6s][WSAPICLI]", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "HR":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-6s][WSAPICLI]", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.PartsExclude = []string{
			zerolog.TimestampFieldName,
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	default:
		log.Logger = log.With().Logger().Output(os.Stderr)
	}
}
