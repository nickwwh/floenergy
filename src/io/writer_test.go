package io_test

import (
	"bufio"
	"bytes"
	"testing"

	appio "floenergy.com/exercise/src/io"
)

func TestFlushOutputSliceToFile_NilOrEmpty(t *testing.T) {
	// Nil slice pointer: should be a no-op and no output
	{
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		var nilSlice *[]appio.MeterReading
		if err := appio.FlushOutputSliceToFile(nilSlice, w); err != nil {
			t.Fatalf("unexpected error for nil slice: %v", err)
		}
		if err := w.Flush(); err != nil {
			t.Fatalf("flush: %v", err)
		}
		if got := buf.String(); got != "" {
			t.Fatalf("expected no output for nil slice, got %q", got)
		}
	}

	// Empty slice: should be a no-op and no output
	{
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		s := []appio.MeterReading{}
		if err := appio.FlushOutputSliceToFile(&s, w); err != nil {
			t.Fatalf("unexpected error for empty slice: %v", err)
		}
		if err := w.Flush(); err != nil {
			t.Fatalf("flush: %v", err)
		}
		if got := buf.String(); got != "" {
			t.Fatalf("expected no output for empty slice, got %q", got)
		}
	}
}

func TestFlushOutputSliceToFile_WritesAndResets(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	s := []appio.MeterReading{
		{Nmi: "NMI1", Timestamp: "2024-01-01 00:00:00", Reading: 1},
		{Nmi: "NMI1", Timestamp: "2024-01-01 00:30:00", Reading: 2.5},
	}

	if err := appio.FlushOutputSliceToFile(&s, w); err != nil {
		t.Fatalf("FlushOutputSliceToFile returned error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	want := "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES \n" +
		"('NMI1', '2024-01-01 00:00:00', 1),\n" +
		"('NMI1', '2024-01-01 00:30:00', 2.5);\n"

	if got := buf.String(); got != want {
		t.Fatalf("SQL output mismatch\nGOT:\n%q\nWANT:\n%q", got, want)
	}

	if len(s) != 0 {
		t.Fatalf("expected slice to be reset to length 0, got len=%d", len(s))
	}
}

func TestFlushOutputSliceToFile_SingleRowFormat(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	s := []appio.MeterReading{
		{Nmi: "NMI2", Timestamp: "2024-02-01 12:00:00", Reading: 3},
	}

	if err := appio.FlushOutputSliceToFile(&s, w); err != nil {
		t.Fatalf("FlushOutputSliceToFile returned error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	want := "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES \n" +
		"('NMI2', '2024-02-01 12:00:00', 3);\n"

	if got := buf.String(); got != want {
		t.Fatalf("single-row SQL output mismatch\nGOT:\n%q\nWANT:\n%q", got, want)
	}

	if len(s) != 0 {
		t.Fatalf("expected slice to be reset to length 0, got len=%d", len(s))
	}
}
