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

	floio "floenergy.com/exercise/src/io"
	"floenergy.com/exercise/src/nem12"
)

const batchSize = 50

func main() {
	inPath := flag.String("in", "resources/test/data/sample.csv", "Path to input NEM12 CSV inputFile")
	outPath := flag.String("out", "output/output.sql", "Path to write generated SQL inserts")
	flag.Parse()

	inputFile, err := os.Open(*inPath)
	if err != nil {
		log.Fatalf("Failed to open inputFile: %v", err)
	}
	outFile, err := os.Create(*outPath)
	//outFile, err := os.OpenFile(*outPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
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

	currentReadInfo := floio.CurrentReaderInfo{
		LineNum: 0,
	}
	var tempSlice []floio.MeterReading
	for {
		record, err := reader.Read()
		currentReadInfo.LineNum++

		if err != nil {
			if err == io.EOF {
				break
			}
			err := fmt.Errorf("error reading line %d: %w", currentReadInfo.LineNum, err)
			slog.Error(err.Error())
			continue
		}
		if len(record) == 0 {
			slog.Warn(fmt.Sprintf("Empty record at line %d", currentReadInfo.LineNum))
			continue
		}

		recordType := record[0]
		switch recordType {
		case "100":
			if record[1] != "NEM12" {
				log.Fatalf("inputFile %s is not in NEM12 format", inputFile.Name())
			}
		case "200":
			err = nem12.Handle200(record, &currentReadInfo)
			if err != nil {
				log.Fatalf("Error processing record at line %d: %v", currentReadInfo.LineNum, err)
			}
		case "300":
			err := nem12.Handle300(record, currentReadInfo, &tempSlice)
			if err != nil {
				log.Fatalf("Error processing record at line %d: %v", currentReadInfo.LineNum, err)
			}
			if len(tempSlice) >= batchSize {
				err := floio.FlushOutputSliceToFile(&tempSlice, writer)
				if err != nil {
					log.Fatalf("Error writing to outputFile: %v", err)
				}
			}

		case "400":
			// not part of the scope but ideally included to indicate where reading quality is less than ideal
		case "500":
		case "900":
			slog.Debug("End of data")
		}
	}

	// write any remaining records

	err = floio.FlushOutputSliceToFile(&tempSlice, writer)
	if err != nil {
		log.Fatalf("Error writing to outputFile: %v", err)
	}

	slog.Info(fmt.Sprintf("Total lines processed: %d", currentReadInfo.LineNum))
}
