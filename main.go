package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	inPath := flag.String("in", "resources/test/data/sample.csv", "Path to input NEM12 CSV inputFile")
	outPath := flag.String("out", "output/output.sql", "Path to write generated SQL inserts")
	flag.Parse()

	inputFile, err := os.Open(*inPath)
	if err != nil {
		log.Fatalf("Failed to open inputFile: %v", err)
	}
	outFile, err := os.OpenFile(*outPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Fatalf("Failed to create output inputFile: %v", err)
	}
	writer := bufio.NewWriter(outFile)
	defer func(inputFile *os.File, writer *bufio.Writer, outFile *os.File) {
		if err := inputFile.Close(); err != nil {
			log.Fatalf("Failed to close inputFile: %v", err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatalf("Failed to flush output inputFile: %v", err)
		}
		if err := outFile.Close(); err != nil {
			log.Fatalf("Failed to close output inputFile: %v", err)
		}

	}(inputFile, writer, outFile)

	reader := csv.NewReader(inputFile)
	reader.FieldsPerRecord = -1

	lineNum := 0
	nmi := ""
	intervals := 0
	intervalPeriod := 0
	for {
		record, err := reader.Read()
		lineNum++

		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("Error reading line %d: %v", lineNum, err)
			continue
		}
		if len(record) == 0 {
			slog.Warn("Empty record at line %d", lineNum)
			continue
		}

		recordType := record[0]
		switch recordType {
		case "100":
			if record[1] != "NEM12" {
				log.Fatal("inputFile is not in NEM12 format")
			}
		case "200":
			if len(record) < 8 {
				log.Fatalf("Record %d is too short: %d columns", lineNum, len(record))
			}
			nmi = record[1]
			intervalPeriod, err = strconv.Atoi(record[8])
			if err != nil {
				return
			}
			intervals = 1440 / intervalPeriod
		case "300":
			if len(record) < intervals+2 {
				log.Fatalf("Record too short, expected %d values but got %d", intervals+2, len(record))
			}

			date, err := time.Parse("20060102", record[1])
			if err != nil {
				log.Fatalf("Error parsing date %s at line %d: %v", record[1], lineNum, err)
			}

			// iterate through the intervals for that date
			for i := 2; i <= intervals+1; i++ {
				reading := record[i]
				readingFloat, err := strconv.ParseFloat(reading, 64)
				if err != nil {
					log.Fatalf("Error parsing reading value %s at line %d: %v", reading, lineNum, err)
				}
				timestamp := date.Add(time.Duration((i-1)*intervalPeriod) * time.Minute)
				_, err = fmt.Fprintf(writer, "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES ('%s', '%s', %v);\n",
					nmi,
					timestamp.Format("2006-01-02 15:04:05"),
					readingFloat)
				if err != nil {
					log.Fatalf("Error writing to inputFile: %v", err)
				}
			}
		case "400":
			// not part of the scope but ideally included to indicate where reading quality is less than ideal
		case "500":
		case "900":
			slog.Debug("End of data")
		}
	}

	slog.Info(fmt.Sprintf("Total lines processed: %d", lineNum))
}
