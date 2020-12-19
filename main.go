package main

import (
	"flag"
	"github.com/klyngen/IGC-parser/parser"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var logRecords = parser.FlightRecordBase{}

func main() {
	// Parse command-line flags

	// Should be able to select a directory
	var directory string
	var summary bool
	flag.StringVar(&directory, "dir", ".", "Give relative path to a directory")
	flag.BoolVar(&summary, "summary", false, "Creates day-based summaries for logs")

	flag.Parse()

	walkDirectories(directory)

	if summary {
		logRecords.PrintStats()
	} else {
		logRecords.PrintAllRecords()
	}
}



func walkDirectories(dir string) {
	content, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Printf("Unable to read directory due to unexpected error: %s", err)
		return
	}

	for _, c := range content {
		// If it is a directory
		// Recurse this motherFucker
		if c.IsDir() {
			newPath := path.Join(dir, c.Name())
			walkDirectories(newPath)
			continue // We dont want to try to parse this dir
		}

		// This is a logfile. This should be parsed
		if strings.HasSuffix(strings.ToUpper(c.Name()), ".IGC") {
			rec, err := parser.ParseFile(path.Join(dir, c.Name()))
			if err == nil {
				// Has to be at least 2 minutes
				if rec.Duration > 120000 {
					logRecords.AddRecord(rec)
				}
			}
		}
	}
}
