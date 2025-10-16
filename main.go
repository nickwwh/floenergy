package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.Open("sample.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	outFile, err := os.Create("output.sql")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	lineNum := 0
	nmi := ""
	intervals := 0
	intervalPeriod := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading line %d: %v", lineNum, err)
		}

		lineNum++

		recordType := record[0]

		// TODO: handle if there are missing/unexpected data all over
		switch recordType {
		case "100":
			if record[1] != "NEM12" {
				err := fmt.Errorf("file is not in NEM12 format")
				fmt.Println(err)
				return
			}
		case "200":
			nmi = record[1]
			intervalPeriod, err = strconv.Atoi(record[8])
			if err != nil {
				return
			}
			intervals = 1440 / intervalPeriod
		case "300":
			date, err := time.Parse("20060102", record[1])
			if err != nil {
				return
			}

			for i := 2; i < intervals+1; i++ {
				reading := record[i]

				readingFloat, _ := strconv.ParseFloat(reading, 64)
				timestamp := date.Add(time.Duration((i-1)*intervalPeriod) * time.Minute)
				fmt.Fprintf(writer, "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES ('%s', '%s', %v);\n",
					nmi,
					timestamp.Format("2006-01-02 15:04:05"),
					readingFloat)
			}
		case "400":
		case "500":
		case "900":
			fmt.Println("End of data")
		}
	}

	fmt.Printf("\nTotal lines processed: %d\n", lineNum)
}
