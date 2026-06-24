package filters

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestFilterToTime(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		want string
	}{
		{name: "under minute", in: 59, want: "00:59"},
		{name: "minute and second", in: 61, want: "01:01"},
		{name: "over an hour", in: 3661, want: "01:01:01"},
		{name: "exactly one hour keeps mm:ss format", in: 3600, want: "60:00"},
		{name: "null float valid", in: sql.NullFloat64{Float64: 90, Valid: true}, want: "01:30"},
		{name: "null float invalid", in: sql.NullFloat64{Valid: false}, want: "00:00"},
		{name: "pgx int8 valid", in: pgtype.Int8{Int64: 125, Valid: true}, want: "02:05"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterToTime(pongo2.AsValue(tc.in), pongo2.AsValue(nil))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tc.want {
				t.Fatalf("got %q, want %q", got.String(), tc.want)
			}
		})
	}
}

func TestFilterSplit(t *testing.T) {
	got, err := filterSplit(pongo2.AsValue(" a|b|c "), pongo2.AsValue("|"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"a", "b", "c"}
	actual, ok := got.Interface().([]string)
	if !ok {
		t.Fatalf("unexpected type %T", got.Interface())
	}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("got %#v, want %#v", actual, want)
	}
}

func TestFilterPrintDate(t *testing.T) {
	tm := time.Date(2025, time.January, 2, 15, 4, 5, 0, time.UTC)

	t.Run("sql null time valid", func(t *testing.T) {
		got, err := filterPrintDate(pongo2.AsValue(sql.NullTime{Time: tm, Valid: true}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "02-01-2025" {
			t.Fatalf("got %q, want %q", got.String(), "02-01-2025")
		}
	})

	t.Run("pgx date valid", func(t *testing.T) {
		got, err := filterPrintDate(pongo2.AsValue(pgtype.Date{Time: tm, Valid: true}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "02-01-2025" {
			t.Fatalf("got %q, want %q", got.String(), "02-01-2025")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		got, err := filterPrintDate(pongo2.AsValue(sql.NullTime{Valid: false}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "" {
			t.Fatalf("got %q, want empty", got.String())
		}
	})
}

func TestFilterPrintDateHuman(t *testing.T) {
	tm := time.Date(2025, time.January, 2, 15, 4, 5, 0, time.UTC)

	t.Run("sql null time valid", func(t *testing.T) {
		got, err := filterPrintDateHuman(pongo2.AsValue(sql.NullTime{Time: tm, Valid: true}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "January 2, 2025" {
			t.Fatalf("got %q, want %q", got.String(), "January 2, 2025")
		}
	})

	t.Run("pgx date valid", func(t *testing.T) {
		got, err := filterPrintDateHuman(pongo2.AsValue(pgtype.Date{Time: tm, Valid: true}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "January 2, 2025" {
			t.Fatalf("got %q, want %q", got.String(), "January 2, 2025")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		got, err := filterPrintDateHuman(pongo2.AsValue(sql.NullTime{Valid: false}), pongo2.AsValue(nil))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.String() != "" {
			t.Fatalf("got %q, want empty", got.String())
		}
	})
}

func TestFilterStoM(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		want string
	}{
		{name: "seconds only", in: float64(59), want: "00:59"},
		{name: "minute and second", in: float32(61), want: "01:01"},
		{name: "over an hour", in: float64(3661), want: "01:01:01"},
		{name: "pgx float4", in: pgtype.Float4{Float32: 125, Valid: true}, want: "02:05"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterStoM(pongo2.AsValue(tc.in), pongo2.AsValue(nil))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tc.want {
				t.Fatalf("got %q, want %q", got.String(), tc.want)
			}
		})
	}
}

func TestFilterSanitizeFile(t *testing.T) {
	got, err := filterSanitizeFile(pongo2.AsValue("path/to/file.txt"), pongo2.AsValue(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.String() != "path_to_file.txt" {
		t.Fatalf("got %q, want %q", got.String(), "path_to_file.txt")
	}
}

func TestFilterFileSize(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		want string
	}{
		{name: "bytes", in: int(500), want: "500.0 B"},
		{name: "kilobytes", in: int64(1024), want: "1.0 KB"},
		{name: "megabytes", in: int64(1024 * 1024), want: "1.0 MB"},
		{name: "gigabytes", in: int64(1024 * 1024 * 1024), want: "1.0 GB"},
		{name: "invalid null", in: sql.NullInt64{Valid: false}, want: ""},
		{name: "pgx int4", in: pgtype.Int4{Int32: 2048, Valid: true}, want: "2.0 KB"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterFileSize(pongo2.AsValue(tc.in), pongo2.AsValue(nil))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tc.want {
				t.Fatalf("got %q, want %q", got.String(), tc.want)
			}
		})
	}
}
