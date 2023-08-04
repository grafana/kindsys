package kindsys

import (
	"bytes"
	"encoding/json"
)

type Maturity int8

const (
	MaturityUnknown Maturity = iota
	MaturityMerged
	MaturityExperimental
	MaturityStable
	MaturityMature
)

func (s Maturity) String() string {
	switch s {
	case MaturityMerged:
		return "merged"
	case MaturityExperimental:
		return "experimental"
	case MaturityStable:
		return "stable"
	case MaturityMature:
		return "mature"
	}
	return "unknown"
}

// MarshalJSON marshals the enum as a quoted json string
func (s Maturity) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON convert a quoted json string to the enum value
func (s *Maturity) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch j {
	case "merged":
		*s = MaturityMerged
	case "experimental":
		*s = MaturityExperimental
	case "stable":
		*s = MaturityStable
	case "mature":
		*s = MaturityMature
	default:
		*s = MaturityUnknown
	}
	return nil
}
