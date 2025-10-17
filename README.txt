#

General idea

* Stream csv line by line

* specific to nem12

* 200 indicates new nmi
* 500 indicates end of data for nmi
* each 300 lines is probably for a day
* no mention of 400 records but it is mainly to indicate the reason behind less than perfect meter readings in the 300 records


# how to use
`go run main.go -i input.csv -o output.json`

If output.json exists, program will throw an error.

* Possible improvements

* indicate when reading quality is affected (handling 400 records)
* save the unit provided (kWh, kvarh, etc)
* could be more explicit with the interval instead of a single timestamp
* implement nem13 parsing
* logic is currently very specific to nem12 and simple enough to be within a single method. Consider refactoring as complexity increases


# tech used

bufio. Ensures minimal footprint when reading large csv
input args allows running in parallel or even as a docker container