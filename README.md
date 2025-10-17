# General idea

1. Stream and parse csv line by line to reduce memory footprint
2. Parse according to nem12 format
3. Batched writes into output SQL file
4. Multi-row inserts for compatibility with most relational databases

* specific to nem12
  * 200 indicates new nmi
  * 500 indicates end of data for nmi
  * each 300 lines is probably for a day
  * no mention of 400 records but it is mainly to indicate the reason behind less than perfect meter readings in the 300 records


# How to run
1. Clone the repo
2. run `go run src/main.go -in resources/test/data/sample.csv -out output/output1.csv`

If you have the binary, you can run it like this:
`./main -in resources/test/data/sample.csv -out output/output3.sql`

If output.csv exists, the program will throw an error.

# Tech used
1. Go. Lightweight, compiled, fast, and easy to use. Compiles to multiple platforms. Standard library contains all the basic functionality required. Go is supported on Lambdas as well so it is great for serverless applications.


# Possible improvements

* indicate when reading quality is affected (handling 400 records)
* save the unit provided (kWh, etc)
* could be more explicit with the interval instead of a single timestamp
* implement nem13 parsing
* logic is currently very specific to nem12 and simple enough to be within a single method. Consider refactoring as complexity increases


# Design decisions
1. bufio library. Ensures minimal footprint when reading large CSVs.
2. Multirow inserts to ensure compatibility with most relational databases.
3. Input args allows efficient piping to be done in CLI (for processing multiple files).
4. Seperate methods for handling different types of records to ensure readability and easier testing.
5. 