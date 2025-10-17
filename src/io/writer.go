package io

import (
	"bufio"
	"fmt"
)

type MeterReading struct {
	Nmi       string
	Timestamp string
	Reading   float64
}

func FlushOutputSliceToFile(tempSlice *[]MeterReading, writer *bufio.Writer) error {
	if tempSlice == nil || len(*tempSlice) == 0 {
		return nil
	}

	_, err := fmt.Fprintf(writer, "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES \n")
	if err != nil {
		return err
	}

	for i, reading := range *tempSlice {
		_, err = fmt.Fprintf(writer, "('%s', '%s', %v)",
			reading.Nmi,
			reading.Timestamp,
			reading.Reading)
		if err != nil {
			return err
		}

		if i < len(*tempSlice)-1 {
			_, err = fmt.Fprintf(writer, ",\n")
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintf(writer, ";\n")
			if err != nil {
				return err
			}
		}
	}
	*tempSlice = (*tempSlice)[:0]
	return nil
}
