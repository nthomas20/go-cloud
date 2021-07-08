/*
 * Filename: logs.go
 * Author: Nathaniel Thomas
 */

package logs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

// Write : Write a log entry
func Write(logger *log.Logger, o interface{}, debugOptional ...bool) {
	var (
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	debug := false

	if len(debugOptional) > 0 {
		debug = debugOptional[0]
	}

	// Marshal the structure to JSON
	if logRow, err := json.Marshal(o); err == nil {
		logger.Println(string(logRow[:]))

		if debug == true {
			fmt.Println(string(logRow[:]))
		}
	}
}

func logInt64Conv(s string, n string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		fmt.Println("Invalid integer conversion [", n, "]:", s)
		os.Exit(1)
	}

	return i
}

// StartLogging : Start up and return a logger
func StartLogging(n string, filename string, maxFilesString string, maxBytesString string, flagOptional ...int) *log.Logger {
	var (
		flags = 0
	)

	maxFiles := logInt64Conv(maxFilesString, n+"_LOG_ROTATE_MAX_FILES")
	maxBytes := logInt64Conv(maxBytesString, n+"_LOG_ROTATE_BYTES")

	// We'll take an optional log flags parameter
	if len(flagOptional) > 0 && flagOptional[0] != 0 {
		flags = flagOptional[0]
	}

	logger, err := configureWriter(filename, int(maxFiles), maxBytes, "", flags)

	if err != nil {
		fmt.Println("ERROR WITH ", n, " LOGGING:", err)
		os.Exit(1)
	}

	return logger
}
