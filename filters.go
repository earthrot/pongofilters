package filters

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/jackc/pgx/v5/pgtype"
)

func init() {
	pongo2.RegisterFilter("totime", filterToTime)
	pongo2.RegisterFilter("splitlist", filterSplit)
	pongo2.RegisterFilter("printdate", filterPrintDate)
	pongo2.RegisterFilter("printdatehuman", filterPrintDateHuman)
	pongo2.RegisterFilter("stom", filterStoM)
	pongo2.RegisterFilter("sanitizefile", filterSanitizeFile)
	pongo2.RegisterFilter("filesize", filterFileSize)
}

func filterToTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var seconds int = 0
	switch v := in.Interface().(type) {
	case float64:
		seconds = int(v)
	case sql.NullFloat64:
		if v.Valid {
			seconds = int(v.Float64)
		}
	case pgtype.Float8:
		if v.Valid {
			seconds = int(v.Float64)
		}
	case int:
		seconds = v
	case int64:
		seconds = int(v)
	case sql.NullInt64:
		if v.Valid {
			seconds = int(v.Int64)
		}
	case pgtype.Int8:
		if v.Valid {
			seconds = int(v.Int64)
		}
	}

	val := fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
	if seconds > 3600 {
		val = fmt.Sprintf("%02d:%s", seconds/3600, val)
	}

	return pongo2.AsSafeValue(val), nil
}

func filterSplit(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var val string = ""
	var delim string = ","

	switch p := param.Interface().(type) {
	case string:
		delim = p
	}

	switch v := in.Interface().(type) {
	case string:
		val = v
	}

	return pongo2.AsSafeValue(strings.Split(strings.TrimSpace(val), delim)), nil
}

func filterPrintDate(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	switch v := in.Interface().(type) {
	case sql.NullTime:
		if v.Valid {
			return pongo2.AsSafeValue(v.Time.Format("02-01-2006")), nil
		}
	case pgtype.Date:
		if v.Valid {
			return pongo2.AsSafeValue(v.Time.Format("02-01-2006")), nil
		}
	}
	return pongo2.AsSafeValue(""), nil
}

func filterPrintDateHuman(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	switch v := in.Interface().(type) {
	case sql.NullTime:
		if v.Valid {
			return pongo2.AsSafeValue(v.Time.Format("January 2, 2006")), nil
		}
	case pgtype.Date:
		if v.Valid {
			return pongo2.AsSafeValue(v.Time.Format("January 2, 2006")), nil
		}
	}
	return pongo2.AsSafeValue(""), nil
}

func filterStoM(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var seconds int = 0
	//f("%T\n", in.Interface())
	switch v := in.Interface().(type) {
	case sql.NullFloat64:
		if v.Valid {
			seconds = int(v.Float64)
		}
	case pgtype.Float4:
		if v.Valid {
			seconds = int(v.Float32)
		}
	case pgtype.Float8:
		if v.Valid {
			seconds = int(v.Float64)
		}
	case float64:
		seconds = int(v)
	case float32:
		seconds = int(v)
	}

	hours := seconds / (60 * 60)
	seconds -= hours * 60 * 60
	minutes := seconds / 60
	seconds -= minutes * 60

	val := fmt.Sprintf("%02d:%02d", minutes, seconds)
	if hours > 0 {
		val = fmt.Sprintf("%02d:%s", hours, val)
	}

	return pongo2.AsSafeValue(val), nil
}

func filterSanitizeFile(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := ""
	switch v := in.Interface().(type) {
	case string:
		val = v
	}

	val = strings.Replace(val, "/", "_", -1)

	return pongo2.AsSafeValue(val), nil
}

func filterFileSize(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := ""
	var intval int64 = -1
	switch v := in.Interface().(type) {
	case int:
		intval = int64(v)
	case int16:
		intval = int64(v)
	case int32:
		intval = int64(v)
	case int64:
		intval = v
	case sql.NullInt16:
		if v.Valid {
			intval = int64(v.Int16)
		}
	case sql.NullInt32:
		if v.Valid {
			intval = int64(v.Int32)
		}
	case sql.NullInt64:
		if v.Valid {
			intval = v.Int64
		}
	case pgtype.Int2:
		if v.Valid {
			intval = int64(v.Int16)
		}
	case pgtype.Int4:
		if v.Valid {
			intval = int64(v.Int32)
		}
	case pgtype.Int8:
		if v.Valid {
			intval = v.Int64
		}
	}

	if intval >= 0 {
		unit := "B"
		size := float64(intval)
		if size >= 1024 {
			size /= 1024
			unit = "KB"
		}
		if size >= 1024 {
			size /= 1024
			unit = "MB"
		}
		if size >= 1024 {
			size /= 1024
			unit = "GB"
		}
		val = fmt.Sprintf("%.1f %s", size, unit)
	}

	return pongo2.AsSafeValue(val), nil
}
