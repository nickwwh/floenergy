package nem12_test

import (
	"testing"

	io2 "floenergy.com/exercise/src/io"
	"floenergy.com/exercise/src/nem12"
)

func TestHandle200_Valid(t *testing.T) {
	// record fields: index 1 = NMI, index 8 = interval period (minutes)
	record := []string{"200", "NMI1234567", "", "", "", "", "", "", "30"}
	ri := io2.CurrentReaderInfo{}

	if err := nem12.Handle200(record, &ri); err != nil {
		t.Fatalf("Handle200 returned unexpected error: %v", err)
	}

	if ri.Nmi != "NMI1234567" {
		t.Fatalf("Nmi mismatch: got %q", ri.Nmi)
	}
	if ri.IntervalPeriod != 30 {
		t.Fatalf("IntervalPeriod mismatch: got %d", ri.IntervalPeriod)
	}
	if ri.Intervals != 48 { // 1440 / 30
		t.Fatalf("Intervals mismatch: got %d", ri.Intervals)
	}
}

func TestHandle200_ShortRecord(t *testing.T) {
	record := []string{"200", "NMI123"} // too short
	ri := io2.CurrentReaderInfo{}

	if err := nem12.Handle200(record, &ri); err == nil {
		t.Fatalf("expected error for short record, got nil")
	}
}

func TestHandle200_BadPeriod(t *testing.T) {
	record := []string{"200", "NMI0001", "", "", "", "", "", "", "bad"}
	ri := io2.CurrentReaderInfo{}

	if err := nem12.Handle200(record, &ri); err == nil {
		t.Fatalf("expected error for bad period, got nil")
	}
	// Nmi is set before parsing the period; verify it was updated
	if ri.Nmi != "NMI0001" {
		t.Fatalf("Nmi should be set before error; got %q", ri.Nmi)
	}
}

func TestHandle300_ValidIntervals(t *testing.T) {
	// CurrentReaderInfo configured for 3 intervals of 30 minutes
	ri := io2.CurrentReaderInfo{Nmi: "NMI123", Intervals: 3, IntervalPeriod: 30}
	// record: ["300", date, v1, v2, v3]
	record := []string{"300", "20240101", "1.0", "2", "3.5"}
	var tmp []io2.MeterReading

	if err := nem12.Handle300(record, ri, &tmp); err != nil {
		t.Fatalf("Handle300 returned error: %v", err)
	}

	if len(tmp) != 3 {
		t.Fatalf("expected 3 readings appended, got %d", len(tmp))
	}

	// Expect timestamps at 00:30, 01:00, 01:30 based on current logic
	want := []struct {
		ts  string
		val float64
	}{
		{"2024-01-01 00:30:00", 1.0},
		{"2024-01-01 01:00:00", 2.0},
		{"2024-01-01 01:30:00", 3.5},
	}
	for i := range want {
		if tmp[i].Nmi != "NMI123" {
			t.Fatalf("reading %d Nmi mismatch: got %q", i, tmp[i].Nmi)
		}
		if tmp[i].Timestamp != want[i].ts {
			t.Fatalf("reading %d timestamp mismatch: got %q want %q", i, tmp[i].Timestamp, want[i].ts)
		}
		if tmp[i].Reading != want[i].val {
			t.Fatalf("reading %d value mismatch: got %v want %v", i, tmp[i].Reading, want[i].val)
		}
	}
}

func TestHandle300_ShortRecord(t *testing.T) {
	ri := io2.CurrentReaderInfo{Nmi: "NMI123", Intervals: 4, IntervalPeriod: 30}
	// only 2 interval values -> len(record) must be >= intervals + 2
	record := []string{"300", "20240101", "1.0", "2.0"}
	var tmp []io2.MeterReading

	if err := nem12.Handle300(record, ri, &tmp); err == nil {
		t.Fatalf("expected error for short record, got nil")
	}
}

func TestHandle300_BadDate(t *testing.T) {
	ri := io2.CurrentReaderInfo{Nmi: "NMI123", Intervals: 1, IntervalPeriod: 30}
	record := []string{"300", "bad-date", "1.0"}
	var tmp []io2.MeterReading

	if err := nem12.Handle300(record, ri, &tmp); err == nil {
		t.Fatalf("expected error for bad date, got nil")
	}
}

func TestHandle300_BadReadingValue(t *testing.T) {
	ri := io2.CurrentReaderInfo{Nmi: "NMI123", Intervals: 2, IntervalPeriod: 15}
	record := []string{"300", "20240101", "1.0", "bad"}
	var tmp []io2.MeterReading

	if err := nem12.Handle300(record, ri, &tmp); err == nil {
		t.Fatalf("expected error for bad reading value, got nil")
	}
}
