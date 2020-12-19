package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadingFiles(t *testing.T) {
	record, err := ParseFile("test-data.IGC")

	if err != nil {
		t.Failed()
	}

	assert.Equal(t, 546571, record.Duration)
	assert.Equal(t, 744, record.MaxAltitude)
	assert.Equal(t, 233, record.MinAltitude)
}

func TestStatsCreation(t *testing.T) {
	base := FlightRecordBase{}

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now(),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now(),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now(),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now().AddDate(0, -1, 0),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now().AddDate(0, -1, 0),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.AddRecord(FlightRecord{
		Duration:    100000,
		StartTime:   time.Now().AddDate(0, -1, 0),
		MaxAltitude: 0,
		MinAltitude: 0,
	})

	base.PrintStats()



}