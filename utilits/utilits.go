package utilits

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	CustomFloat64 struct {
		Float64 float64
	}
	CustomInt64 struct {
		Int64 int64
	}
	CustomInt struct {
		Int int
	}
	CustomBool bool
	CustomTime time.Time
)

const QUOTES_BYTE = 34

func (cb *CustomBool) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"TRUE"`, `TRUE`, `"true"`, `true`, `"1"`, `1`:
		*cb = true
	case `"FALSE"`, `FALSE`, `"false"`, `false`, `"0"`, `0`, `""`:
		*cb = false
	default:
		return fmt.Errorf(`CustomBool: parsing "%s": unknown value`, string(data))
	}
	return nil
}

func (cf *CustomFloat64) UnmarshalJSON(data []byte) error {
	if data[0] == QUOTES_BYTE {
		if err := json.Unmarshal(bytes.Replace(data[1:len(data)-1], []byte{44}, []byte{}, -1), &cf.Float64); err != nil {
			return fmt.Errorf("CustomFloat64: UnmarshalJSON with quotes data [%s]: %w", string(data[1:len(data)-1]), err)
		}
	} else {
		if err := json.Unmarshal(data, &cf.Float64); err != nil {
			return fmt.Errorf("CustomFloat64: UnmarshalJSON without quotes data [%s]: %w", string(data), err)
		}
	}
	return nil
}

func (ci *CustomInt64) UnmarshalJSON(data []byte) error {
	if data[0] == QUOTES_BYTE {
		if string(data[1:len(data)-1]) == "-" {
			ci.Int64 = 0
			return nil
		}
		if err := json.Unmarshal(data[1:len(data)-1], &ci.Int64); err != nil {
			return fmt.Errorf("CustomInt64: UnmarshalJSON: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &ci.Int64); err != nil {
			return fmt.Errorf("CustomInt64: UnmarshalJSON: %w", err)
		}
	}
	return nil
}

func (ci *CustomInt) UnmarshalJSON(data []byte) error {
	if data[0] == QUOTES_BYTE {
		if string(data[1:len(data)-1]) == "-" {
			ci.Int = 0
			return nil
		}
		if err := json.Unmarshal(data[1:len(data)-1], &ci.Int); err != nil {
			return fmt.Errorf("CustomInt: UnmarshalJSON: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &ci.Int); err != nil {
			return fmt.Errorf("CustomInt: UnmarshalJSON: %w", err)
		}
	}
	return nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if data[0] == QUOTES_BYTE {
		if string(data[1:len(data)-1]) == "-" {
			*ct = CustomTime{}
			return nil
		}

		t, err := time.Parse("2006-01-02 15:04:05", string(data[1:len(data)-1]))
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05.000000", string(data[1:len(data)-1]))
			if err != nil {
				*ct = CustomTime{}
				return nil
			} else {
				*ct = CustomTime(t)
				return nil
			}
		}
		*ct = CustomTime(t)
		return nil
	}
	return fmt.Errorf(`CustomTime: parsing "%s": unknown value`, string(data))
}

func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func ClearString(s string) string {
	rep := strings.NewReplacer(
		"~", "",
		"`", "",
		"!", "",
		"@", "",
		"#", "",
		"$", "",
		"%", "",
		"^", "",
		"&", "",
		"*", "",
		"(", "",
		")", "",
		"_", "",
		"+", "",
		"-", "",
		"=", "",
		"{", "",
		"}", "",
		"[", "",
		"]", "",
		",", "",
		"/", "",
		"?", "",
		":", "",
		"<", "",
		">", "",
		"'", "",
		";", "",
		`\`, "",
		`"`, "",
		"|", "",
		"№", "",
		" ", "",
	)
	return rep.Replace(s)
}
