package nem12

import (
	"errors"
	"strconv"
	"time"

	"floenergy.com/exercise/src/io"
)

func Handle300(record []string, currentReadInfo io.CurrentReaderInfo, tempSlice *[]io.MeterReading) (err error) {
	if len(record) < currentReadInfo.Intervals+2 {
		return errors.New("record shorter than expected")
	}

	date, err := time.Parse("20060102", record[1])
	if err != nil {
		return err
	}

	// iterate through the intervals for that date
	for i := 2; i <= currentReadInfo.Intervals+1; i++ {
		reading := record[i]
		readingFloat, err := strconv.ParseFloat(reading, 64)
		if err != nil {
			return err
		}
		timestamp := date.Add(time.Duration((i-1)*currentReadInfo.IntervalPeriod) * time.Minute)

		*tempSlice = append(*tempSlice, io.MeterReading{currentReadInfo.Nmi, timestamp.Format("2006-01-02 15:04:05"), readingFloat})
	}
	return nil
}

func Handle200(record []string, currentReadInfo *io.CurrentReaderInfo) error {
	if len(record) < 8 {
		return errors.New("record shorter than expected")
	}

	var err error
	currentReadInfo.Nmi = record[1]
	currentReadInfo.IntervalPeriod, err = strconv.Atoi(record[8])
	if err != nil {
		return err
	}
	currentReadInfo.Intervals = 1440 / currentReadInfo.IntervalPeriod
	return nil
}
