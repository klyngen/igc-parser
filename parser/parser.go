package parser

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// FlightRecord describes one single flight
type FlightRecord struct {
	Duration    int
	StartTime   time.Time
	MaxAltitude int
	MinAltitude int
}

type FlightRecordDay struct {
	records []*FlightRecord
	date time.Time
}

func (f * FlightRecordDay) IsTheDay(date time.Time) bool {
	if f.date.Day() == date.Day() && f.date.Month() == date.Month() {
		return true
	}
	return false
}

func (f *FlightRecordDay) PrintRecords() {
	for _, r := range f.records {
		r.Print()
	}
}

type FlightRecordBase struct {
	recordDays []*FlightRecordDay
}

func (f * FlightRecordBase) AddRecord(record * FlightRecord) {
	if f.recordDays == nil {
		f.recordDays = make([]*FlightRecordDay, 0)
	}
	for _, day := range f.recordDays {
		if day.IsTheDay(record.StartTime) {
			day.Add(record)
			return
		}
	}
	day := FlightRecordDay{
		records: make([]*FlightRecord, 0),
		date: record.StartTime,
	}

	day.Add(record)

	f.recordDays = append(f.recordDays, &day)
}

func (f *FlightRecordDay) Add(record * FlightRecord) {
	f.records = append(f.records, record)
}

func (f *FlightRecordBase) PrintAllRecords() {
	for _, d := range f.recordDays {
		d.PrintRecords()
	}
}

func (f *FlightRecordBase) PrintStats() {
	for _, d := range f.recordDays {
		d.PrintDayStats()
	}
}

func (f * FlightRecordDay) PrintDayStats() {
	var duration = 0.0
	for _, r := range f.records {
		duration += r.GetSeconds()
	}

	fmt.Printf(
`
	DATE: %v
		Number of flights: %v
		FlightTime: 	   h%v:m%v:s%v

`, f.date, len(f.records),
		math.Floor(duration/60.0/60.0),
		math.Floor(math.Mod(duration/60.0, 60.0)),
		math.Floor(math.Mod(duration/60, 1)*60))
}

func (rec * FlightRecord) Print() {

	seconds := rec.GetSeconds()
	fmt.Printf(`

	FLIGHT: %v
	============
	MaxHeight: 	%v
	MinHeight: 	%v
	Duration: 	h%v:m%v:s%v

`, rec.StartTime,
rec.MaxAltitude,
rec.MinAltitude,
math.Floor(seconds/60.0/60.0),
math.Floor(math.Mod(seconds/60.0, 60.0)),
math.Floor(math.Mod(seconds/60, 1)*60))
}

func (rec * FlightRecord) GetSeconds() float64 {
	toDateTime := rec.StartTime
	toDateTime = toDateTime.Add(time.Millisecond * time.Duration(rec.Duration))

	diff := toDateTime.Sub(rec.StartTime)
	return diff.Seconds()
}

func ParseFile(fileName string) (*FlightRecord, error) {
	bytes, err := os.OpenFile(fileName, os.O_RDONLY, 0755)

	defer bytes.Close()

	// Something is wrong
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Seek to the end
	bytes.Seek(0, 2)

	// First line is bullshit
	for i := 0; i < 4; i++ {
		getLine(bytes) // We dont really care about these lines
	}

	record := FlightRecord{
		MinAltitude: getIntegerValue(getLine(bytes)),
		MaxAltitude: getIntegerValue(getLine(bytes)),
		Duration: getIntegerValue(getLine(bytes)),
		StartTime: time.Unix(int64(getIntegerValue(getLine(bytes))), 0),
	}

	// We probably have something
	return &record, nil

}

func getIntegerValue(value string) int {
	// Looks like the following:
	// LXSB  SKYDROP-ALT-MIN-m: 233
	val := strings.Split(value, ": ")[1]

	integer, err := strconv.Atoi(strings.Trim(val, " "))

	if err == nil {
		return integer
	}
	return 0
}

// This piece of code lets us read a file backwards. Apparently... Noone does this....
func getLine(file *os.File) string {
	var bytes []byte = make([]byte, 32)

	// First get the current position
	var firstBreak int64 = 0

	// Start looking for the \n
	for {
		// Seek backwards
		file.Seek(-32, 1)
		// Read moves the pointer
		file.Read(bytes)
		// Move the pointer back
		file.Seek(-32, 1)

		// Check if we have a line break here
		for i := len(bytes)-1; i >= 0; i-- {
			b := bytes[i]
			if b == 13 { // We have a break

				// Get the current pos
				currentPos, _ := file.Seek(0, 1)
				breakLocation := currentPos + int64(i)

				if firstBreak == 0 {
					firstBreak = breakLocation
				} else {
					sz := firstBreak - breakLocation

					res := make([]byte, sz)

					file.Seek(breakLocation, 0)
					file.Read(res)

					// Seek back again
					file.Seek(breakLocation+1, 0)

					return string(res)
				}


			}
		}
	}

}
